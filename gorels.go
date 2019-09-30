package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kopoli/gorels/options"
)

var (
	version     = "Undefined"
	timestamp   = "Undefined"
	buildGOOS   = "Undefined"
	buildGOARCH = "Undefined"
	progVersion = "" + version
)

func fault(err error, message string, arg ...string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s%s: %s\n", message, strings.Join(arg, " "), err)
		os.Exit(1)
	}
}

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

func newCmdMap() *cmdMap {
	var ret cmdMap = make(cmdMap)

	ret.add("bump-major", "Bump the major version number", func(s string) {
	})
	ret.add("bump-minor", "Bump the minor version number", func(s string) {
	})
	ret.add("bump-patch", "Bump the patch level version number", func(s string) {
	})
	ret.add("set-version=", "Bump the patch level version number", func(s string) {
	})

	return &ret
}

type versionData struct {
	
}


func main() {

	opts := options.New()

	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	optVersion := fs.Bool("version", false, "Display version")
	optDebug := fs.Bool("debug", false, "Enable debug output")
	optDryRun := fs.Bool("dryrun", false, "Don't actually run any commands")
	optMessage := fs.String("message", "", "String to add to the tag object")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: Tag commits with semantic versions\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Command line options:")
		fs.PrintDefaults()
	}

	fs.Parse(os.Args[1:])

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

	opts.Set("message", *optMessage)

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}
}
