package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/kopoli/gorels/options"
)

var (
	version     = "Undefined"
	timestamp   = "Undefined"
	buildGOOS   = "Undefined"
	buildGOARCH = "Undefined"
	progVersion = "" + version
)

type cmdDesc struct {
	help string
	op   func(string)
}

type opMap map[string]cmdDesc

func (c *opMap) add(name, help string, op func(string)) {
	(*c)[name] = cmdDesc{
		help: help,
		op:   op,
	}
}

func (c *opMap) help(out io.Writer) {
	names := []string{}
	for k := range *c {
		names = append(names, k)
	}
	sort.Strings(names)

	wr := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(wr, "Commands:")
	printCmd := func(i int) {
		fmt.Fprintf(wr, "  %s\t%s\n", names[i], (*c)[names[i]].help)
	}
	for i := range names {
		printCmd(i)
	}
	wr.Flush()
}

func cmdStr(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func cmdStrOneLine(args ...string) string {
	s, err := cmdStr(args...)
	if err != nil {
		return ""
	}
	return strings.Trim(s, " \n\r\t")
}

type Git struct {
	Git        string
	Commit     string
	TagPrefix  string
	Tags       []string
	DryRun     bool
	debugPrint func(...interface{})
}

func (g *Git) GetTags() error {
	s, err := cmdStr(g.Git, "log", "--format=%D")
	if err != nil {
		return err
	}

	for {
		tagstr := "tag: "
		start := strings.Index(s, tagstr)
		if start < 0 {
			break
		}
		s = s[start+len(tagstr):]
		end := strings.IndexAny(s, ",\n\r")
		if end < 0 {
			break
		}
		tag := s[:end]
		if tag != "" {
			g.Tags = append(g.Tags, tag)
		}
		s = s[end+1:]
	}

	return nil
}

func (g *Git) GetShortLog(start, end string) (string, error) {
	commitrange := start + ".." + end
	if start == "" {
		commitrange = end
	}
	return cmdStr(g.Git, "shortlog", commitrange)
}

func (g *Git) CreateTag(version, message string) error {
	prevtag := ""
	if len(g.Tags) > 0 {
		prevtag = g.Tags[0]
	}
	shortlog, err := g.GetShortLog(prevtag, g.Commit)
	if err != nil {
		return err
	}

	message = fmt.Sprintf("%s\n\n%s", message, shortlog)

	cmdline := []string{g.Git, "tag", "--annotate", "-m", message, version,
		g.Commit}
	g.debugPrint("Running:", strings.Join(cmdline, " "))

	if g.DryRun {
		return nil
	}
	_, err = cmdStr(cmdline...)
	return err
}

type versionData struct {
	operations opMap

	debug   bool
	message string
	version SemVer
	err     error
	git     Git
}

func newVersionData(opts options.Options) *versionData {
	ret := &versionData{
		debug: opts.IsSet("debug"),
		git: Git{
			Git:        "git",
			Commit:     "HEAD",
			TagPrefix:  "v",
			DryRun:     opts.IsSet("dryrun"),
			debugPrint: func(...interface{}) {},
		},
	}
	t := make(opMap)

	debugPrint := func(args ...interface{}) {
		if !ret.debug {
			return
		}
		fmt.Printf(">> ")
		fmt.Println(args...)
	}

	if ret.debug {
		ret.git.debugPrint = debugPrint
	}

	t.add("git=", "Git program to use.", func(s string) {
		debugPrint("Setting git to", s)
		ret.git.Git = s
	})
	t.add("bump-major", "Bump the major version number.", func(s string) {
		debugPrint("Bumping major version")
		ret.version.BumpMajor()
	})
	t.add("bump-minor", "Bump the minor version number.", func(s string) {
		debugPrint("Bumping minor version")
		ret.version.BumpMinor()
	})
	t.add("bump-patch", "Bump the patch level version number.", func(s string) {
		debugPrint("Bumping patch level")
		ret.version.BumpPatch()
	})
	t.add("set-version=", "Set explicit version.", func(s string) {
		debugPrint("Setting version to", s)
		ret.err = ret.version.Set(s)
	})
	t.add("commit=", "Commit to operate on. Default: HEAD", func(s string) {
		debugPrint("Setting git commit to:", s)
		ret.git.Commit = s
	})
	t.add("message=", "Message to inject into the tag", func(s string) {
		debugPrint("Injecting message to tag:", s)
		ret.message = s
	})
	t.add("set-tag-prefix=", "Set tag prefix. Default 'v'.", func(s string) {
		debugPrint("Setting the tag prefix to:", s)
		ret.git.TagPrefix = s
	})
	t.add("tag", "Create a tag.", func(s string) {
		verstr := ret.git.TagPrefix + ret.version.String()
		debugPrint("Creating the git tag:", verstr)
		if ret.message != "" {
			debugPrint("Injecting message:", ret.message)
		}
		ret.git.CreateTag(verstr, ret.message)
		ret.message = ""
	})
	t.add("amend", "Amend the current tag.", func(s string) {
	})
	ret.operations = t

	return ret
}

func parseOp(name string) string {
	return strings.SplitAfter(name, "=")[0]
}

func (v *versionData) checkOperations(operations ...string) error {
	inv := make(map[string]bool)

	for i := range operations {
		n := parseOp(operations[i])
		if _, ok := v.operations[n]; !ok {
			inv[n] = true
		}
	}

	suffix := "s"
	switch len(inv) {
	case 0:
		return nil
	case 1:
		suffix = ""
	}
	var invalid []string
	for k := range inv {
		invalid = append(invalid, k)
	}

	return fmt.Errorf("Invalid operation%s: %s", suffix,
		strings.Join(invalid, ", "))
}

func (v *versionData) apply(operations ...string) error {
	for i := range operations {
		n := parseOp(operations[i])
		if t, ok := v.operations[n]; ok {
			arg := ""
			if strings.Contains(n, "=") {
				arg = strings.SplitN(operations[i], "=", 2)[1]
			}
			t.op(arg)
			if v.err != nil {
				return fmt.Errorf("Operation \"%s\" failed with: %v",
					operations[i], v.err)
			}
		}
	}
	return nil
}

func (v *versionData) getPreviousVersion() error {
	err := v.git.GetTags()
	if err != nil {
		return fmt.Errorf("Getting git tags failed with: %v", err)
	}

	if len(v.git.Tags) == 0 {
		return nil
	}

	prevVersion := v.git.Tags[0]

	v.git.debugPrint("Found", prevVersion, "as previous version")

	err = v.version.Set(prevVersion)
	if err != nil {
		return fmt.Errorf("Parsing previous version \"%s\" failed with: %v",
			prevVersion, err)
	}

	return nil
}

func fault(err error, message string, arg ...string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s%s: %s\n",
			message, strings.Join(arg, " "), err)
		os.Exit(1)
	}
}

func main() {

	opts := options.New()

	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)

	var (
		optVersion = false
		optList = false
		optDebug = false
		optDryRun = false
	)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	opt := func(optvar *bool, longflg, flg, help string) {
		fs.BoolVar(optvar, longflg, false, help)
		if flg != "" {
			fs.BoolVar(optvar, flg, false, help + " (shorthand)")
		}
	}
	opt(&optVersion, "version", "v", "Display version.")
	opt(&optList, "list", "l", "List operations.")
	opt(&optDebug, "debug", "", "Enable debug output.")
	opt(&optDryRun, "dryrun", "d", "Don't actually run any operations. Implies -debug.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: Tag commits with semantic versions\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Command line options:")
		fs.PrintDefaults()
	}

	err := fs.Parse(os.Args[1:])
	fault(err, "Parsing the command line failed")

	if optVersion {
		fmt.Println(options.VersionString(opts))
		os.Exit(0)
	}

	if optDebug {
		opts.Set("debug", "t")
	}

	if optDryRun {
		opts.Set("dryrun", "t")
		opts.Set("debug", "t")
	}

	vd := newVersionData(opts)

	if optList {
		vd.operations.help(os.Stdout)
		os.Exit(0)
	}
	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	err = vd.checkOperations(args...)
	fault(err, "Validating given operations failed")

	err = vd.getPreviousVersion()
	fault(err, "Getting previous version failed")

	err = vd.apply(args...)
	fault(err, "Applying operations failed")
}
