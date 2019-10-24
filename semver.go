package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regular expression for parsing semver from https://semver.org/
	semverRe = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
)

type SemVer struct {
	major      uint64
	minor      uint64
	patch      uint64
	prerelease string
	build      string
}

func (v *SemVer) String() string {
	extra := ""
	if v.prerelease != "" {
		extra = extra + "-" + v.prerelease
	}
	if v.build != "" {
		extra = extra + "+" + v.build
	}
	return fmt.Sprintf("%d.%d.%d%s", v.major, v.minor, v.patch, extra)
}

func (v *SemVer) Set(version string) error {
	version = strings.TrimSpace(version)

	versions := semverRe.FindAllStringSubmatch(version, 1)
	if len(versions) == 0 {
		return fmt.Errorf("Invalid semantic version: %s", version)
	}
	components := versions[0]

	parseUint := func(s string) uint64 {
		ret, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			// This would mean the regexp parsed an improper number
			msg := fmt.Sprintf("Internal error on parsing version: %s", s)
			panic(msg)
		}
		return ret
	}

	v.major = parseUint(components[1])
	v.minor = parseUint(components[2])
	v.patch = parseUint(components[3])
	v.prerelease = components[4]
	v.build = components[5]

	return nil
}

func (v *SemVer) BumpMajor() {
	*v = SemVer{
		major: v.major + 1,
	}
}

func (v *SemVer) BumpMinor() {
	*v = SemVer{
		major: v.major,
		minor: v.minor + 1,
	}
}

func (v *SemVer) BumpPatch() {
	*v = SemVer{
		major: v.major,
		minor: v.minor,
		patch: v.patch + 1,
	}
}

func (v *SemVer) SetPreRelease(prerelease string) {
	v.prerelease = prerelease
	v.build = ""
}

func (v *SemVer) SetBuild(build string) {
	v.build = build
}
