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
	cmd  func(string)
}

type cmdMap map[string]cmdDesc

func (c *cmdMap) add(name, help string, cmd func(string)) {
	(*c)[name] = cmdDesc{
		help: help,
		cmd:  cmd,
	}
}

func (c *cmdMap) help(out io.Writer) {
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
	Git    string
	Commit string
	Tags   []string
	DryRun bool
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

func (g *Git) CreateTag(version string) error {

	return nil
}

type versionData struct {
	commands cmdMap

	debug   bool
	message string
	version SemVer
	err     error
	git     Git
}

func newVersionData(opts options.Options) *versionData {
	ret := &versionData{
		debug:   opts.IsSet("debug"),
		git: Git{
			Git:    "git",
			Commit: "HEAD",
			DryRun: opts.IsSet("dryrun"),
		},
	}
	t := make(cmdMap)

	t.add("git=", "Git program to use.", func(s string) {
		ret.git.Git = s
	})
	t.add("bump-major", "Bump the major version number.", func(s string) {
		ret.err = ret.version.BumpMajor()
	})
	t.add("bump-minor", "Bump the minor version number.", func(s string) {
		ret.err = ret.version.BumpMinor()
	})
	t.add("bump-patch", "Bump the patch level version number.", func(s string) {
		ret.err = ret.version.BumpPatch()
	})
	t.add("set-version=", "Set explicit version.", func(s string) {
		ret.err = ret.version.Set(s)
	})
	t.add("commit=", "Commit to operate on. Default: HEAD", func(s string) {
		ret.git.Commit = s
	})
	t.add("message=", "Message to inject into the tag", func(s string) {
		ret.message = s
	})
	t.add("tag", "Create a tag.", func(s string) {

		ret.message = ""
	})
	t.add("amend", "Amend the current tag.", func(s string) {
	})
	ret.commands = t

	return ret
}

func (v *versionData) checkCommands(commands ...string) error {
	return nil
}

func (v *versionData) apply(commands ...string) error {
	return nil
}

func fault(err error, message string, arg ...string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s%s: %s\n", message, strings.Join(arg, " "), err)
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

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	optVersion := fs.Bool("version", false, "Display version.")
	optList := fs.Bool("list", false, "List commands.")
	optDebug := fs.Bool("debug", false, "Enable debug output.")
	optDryRun := fs.Bool("dryrun", false, "Don't actually run any commands. Implies -debug.")
	optVersionPrefix := fs.String("version-prefix", "v", "String prefix to be stripped when evaluating git version tags.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: Tag commits with semantic versions\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Command line options:")
		fs.PrintDefaults()
	}

	err := fs.Parse(os.Args[1:])
	fault(err, "Parsing the command line failed")

	if *optVersion {
		fmt.Println(options.VersionString(opts))
		os.Exit(0)
	}

	if *optDebug {
		opts.Set("debug", "t")
	}

	if *optDryRun {
		opts.Set("dryrun", "t")
	}

	opts.Set("version-prefix", *optVersionPrefix)

	vd := newVersionData(opts)

	if *optList {
		vd.commands.help(os.Stdout)
		os.Exit(0)
	}
	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	err = vd.checkCommands(args...)
	fault(err, "Validating given commands failed")

	vd.git.GetTags()
	fmt.Println(vd.git.Tags)
}
