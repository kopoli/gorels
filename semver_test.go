package main

import "testing"

func TestSemVerComplete(t *testing.T) {
	var ver *SemVer

	opSet := func(version string) func() {
		return func() {
			ver.Set(version)
		}
	}
	opBumpMaj := func() {
		ver.BumpMajor()
	}
	opBumpMin := func() {
		ver.BumpMinor()
	}
	opBumpPat := func() {
		ver.BumpPatch()
	}
	opSetPre := func(pr string) func() {
		return func() {
			ver.SetPreRelease(pr)
		}
	}
	opSetBld := func(bld string) func() {
		return func() {
			ver.SetBuild(bld)
		}
	}
	opExpect := func(version string) func() {
		return func() {
			if version != ver.String() {
				t.Errorf("Expected version string to be: %s\nGot: %s",
					version, ver.String())
			}
		}
	}
	tests := []struct {
		name string
		ops  []func()
	}{
		{"Empty version should be zero", []func(){opExpect("0.0.0")}},
		{"Set version 1", []func(){opSet("1.0.0"), opExpect("1.0.0")}},
		{"Set version 2", []func(){opSet("1.1.0"), opExpect("1.1.0")}},
		{"Set version 3", []func(){opSet("1.1.1"), opExpect("1.1.1")}},
		{"Set version 4", []func(){opSet("1.1.1-e"), opExpect("1.1.1-e")}},
		{"Set version 5", []func(){opSet("1.1.1+b"), opExpect("1.1.1+b")}},
		{"Set version 6", []func(){opSet("1.1.1-e+b"), opExpect("1.1.1-e+b")}},
		{"Set version 7", []func(){opSet("1.1.1+b-e"), opExpect("1.1.1+b-e")}},
		{"Bump patch", []func(){opSet("1.1.1-e"), opBumpPat, opExpect("1.1.2")}},
		{"Bump minor", []func(){opSet("1.1.1-e"), opBumpMin, opExpect("1.2.0")}},
		{"Bump major", []func(){opSet("1.1.1-e"), opBumpMaj, opExpect("2.0.0")}},
		{"Bump patch empty", []func(){opBumpPat, opExpect("0.0.1")}},
		{"Bump minor empty", []func(){opBumpMin, opExpect("0.1.0")}},
		{"Bump major empty", []func(){opBumpMaj, opExpect("1.0.0")}},
		{"Bump major minor", []func(){opBumpMaj, opBumpMin, opExpect("1.1.0")}},

		{"Set prerelease empty", []func(){opSetPre("abc"), opExpect("0.0.0-abc")}},
		{"Set build empty", []func(){opSetBld("bld"), opExpect("0.0.0+bld")}},
		{"Set prerelease build empty", []func(){opSetPre("pre"), opSetBld("bld"), opExpect("0.0.0-pre+bld")}},
		{"Override prerelease", []func(){opSet("1.0.0-jep"), opSetPre("pre"), opSetBld("bld"), opExpect("1.0.0-pre+bld")}},
		{"Add build", []func(){opSet("1.0.0-jep"), opSetBld("bld"), opExpect("1.0.0-jep+bld")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver = &SemVer{}

			for _, op := range tt.ops {
				op()
			}
		})
	}
}
