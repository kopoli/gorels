package main

import (
	"github.com/blang/semver"
)

// Contractual wrapper around external semantic versioning library

type SemVer struct {
	version semver.Version
}

func (v *SemVer) String() string {
	return v.version.String()
}

func (v *SemVer) Set(version string) error {
	sv, err := semver.ParseTolerant(version)
	if err != nil {
		return err
	}

	v.version = sv
	return nil
}

func (v *SemVer) BumpMajor() error {
	return v.version.IncrementMajor()
}

func (v *SemVer) BumpMinor() error {
	return v.version.IncrementMinor()
}

func (v *SemVer) BumpPatch() error {
	return v.version.IncrementPatch()
}
