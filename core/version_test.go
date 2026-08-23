package core

import "testing"

func TestParseVersion(t *testing.T) {
	cases := []struct {
		in   string
		want Version
		ok   bool
	}{
		{"0.2.0", Version{Major: 0, Minor: 2, Patch: 0}, true},
		{"v0.5.1", Version{Major: 0, Minor: 5, Patch: 1}, true},
		{"1.2.3-alpha+build", Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha", Build: "build"}, true},
		{"", Version{}, false},
		{"bad", Version{}, false},
		{"1.2", Version{}, false},
	}
	for _, c := range cases {
		got, err := ParseVersion(c.in)
		if c.ok && err != nil {
			t.Errorf("ParseVersion(%q) error = %v, want nil", c.in, err)
		}
		if !c.ok && err == nil {
			t.Errorf("ParseVersion(%q) succeeded, want error", c.in)
		}
		if c.ok && got != c.want {
			t.Errorf("ParseVersion(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
}

func TestVersionString(t *testing.T) {
	v := MustParseVersion("0.2.0")
	if s := v.String(); s != "0.2.0" {
		t.Errorf("String() = %q, want 0.2.0", s)
	}
	v2 := MustParseVersion("1.0.0-alpha+001")
	if s := v2.String(); s != "1.0.0-alpha+001" {
		t.Errorf("String() = %q, want 1.0.0-alpha+001", s)
	}
}

func TestVersionCompare(t *testing.T) {
	a := MustParseVersion("0.2.0")
	b := MustParseVersion("0.1.0")
	if a.Compare(b) <= 0 {
		t.Error("0.2.0 should be > 0.1.0")
	}
	if b.Compare(a) >= 0 {
		t.Error("0.1.0 should be < 0.2.0")
	}
	c := MustParseVersion("0.2.0")
	if a.Compare(c) != 0 {
		t.Error("0.2.0 == 0.2.0 should be 0")
	}
	// Pre-release has lower precedence
	pre := MustParseVersion("1.0.0-alpha")
	rel := MustParseVersion("1.0.0")
	if pre.Compare(rel) >= 0 {
		t.Error("1.0.0-alpha should be < 1.0.0")
	}
}

func TestIsCompatible(t *testing.T) {
	// 0.y.z: minor must match
	req := MustParseVersion("0.1.0")
	if !MustParseVersion("0.1.5").IsCompatible(req) {
		t.Error("0.1.5 should be compatible with 0.1.0")
	}
	if MustParseVersion("0.2.0").IsCompatible(req) {
		t.Error("0.2.0 should not be compatible with 0.1.0 (minor mismatch for 0.y.z)")
	}
	// >=1.0.0: major must match
	req1 := MustParseVersion("1.2.0")
	if !MustParseVersion("1.5.0").IsCompatible(req1) {
		t.Error("1.5.0 should be compatible with 1.2.0")
	}
	if MustParseVersion("2.0.0").IsCompatible(req1) {
		t.Error("2.0.0 should not be compatible with 1.2.0")
	}
	if MustParseVersion("1.1.0").IsCompatible(req1) {
		t.Error("1.1.0 should not be compatible with 1.2.0 (actual < required)")
	}
}

func TestCheckEcosystemCompatibility(t *testing.T) {
	required := map[ModuleName]Version{
		ModuleLibs: MustParseVersion("0.1.0"),
	}
	actualOK := map[ModuleName]Version{
		ModuleLibs: MustParseVersion("0.1.5"),
	}
	if err := CheckEcosystemCompatibility(required, actualOK); err != nil {
		t.Errorf("CheckEcosystemCompatibility should succeed, got %v", err)
	}
	actualBad := map[ModuleName]Version{
		ModuleLibs: MustParseVersion("0.2.0"),
	}
	if err := CheckEcosystemCompatibility(required, actualBad); err == nil {
		t.Error("CheckEcosystemCompatibility should fail for 0.2.0 vs 0.1.0")
	}
	actualMissing := map[ModuleName]Version{}
	if err := CheckEcosystemCompatibility(required, actualMissing); err == nil {
		t.Error("should fail for missing module")
	}
}
