package core

import "testing"

func TestParseSpecID(t *testing.T) {
	cases := []struct {
		in string
		ok bool
	}{
		{"KWF-M8K2Q", true},
		{"KWF-ARCH-M8K2Q", true},
		{"KWF-ARCH-M8K2Q-unified-framework-vision", true},
		{"KWL-K1N2Q", true},
		{"KWD-8NQ2K", true},
		{"", false},
		{"KWF-", false},
		{"XXX-M8K2Q", false},
		{"KWF-ARCH-1234", false}, // too short
		{"KWF-ARCH-TOOLONG123", false},
	}
	for _, c := range cases {
		_, err := ParseSpecID(c.in)
		if c.ok && err != nil {
			t.Errorf("ParseSpecID(%q) error = %v, want nil", c.in, err)
		}
		if !c.ok && err == nil {
			t.Errorf("ParseSpecID(%q) succeeded, want error", c.in)
		}
	}
}

func TestSpecIDProjectScopeCode(t *testing.T) {
	id, _ := ParseSpecID("KWF-ARCH-M8K2Q")
	if got := id.Project(); got != "KWF" {
		t.Errorf("Project() = %q, want KWF", got)
	}
	if got := id.Scope(); got != "ARCH" {
		t.Errorf("Scope() = %q, want ARCH", got)
	}
	if got := id.Code(); got != "M8K2Q" {
		t.Errorf("Code() = %q, want M8K2Q", got)
	}
	id2, _ := ParseSpecID("KWF-M8K2Q")
	if got := id2.Scope(); got != "" {
		t.Errorf("Scope() for short form = %q, want empty", got)
	}
	if got := id2.Code(); got != "M8K2Q" {
		t.Errorf("Code() = %q, want M8K2Q", got)
	}
}

func TestParseRequirementID(t *testing.T) {
	if err := ParseRequirementID("FRK-CLI-001"); err != nil {
		t.Errorf("ParseRequirementID valid error = %v", err)
	}
	if err := ParseRequirementID("bad"); err == nil {
		t.Error("ParseRequirementID should fail for bad")
	}
}
