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

func main() {
	optVersion := flag.Bool("v", false, "Display version")

	opts := options.New()

	opts.Set("program-name", os.Args[0])
	opts.Set("program-version", progVersion)
	opts.Set("program-timestamp", timestamp)
	opts.Set("program-buildgoos", buildGOOS)
	opts.Set("program-buildgoarch", buildGOARCH)

	flag.Parse()

	if *optVersion {
		fmt.Println(options.VersionString(opts))
		os.Exit(0)
	}
}
