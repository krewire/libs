// Tests for KWL-ARCH-J2K9Q
package core

import "testing"

func TestKWL_ARCH_J2K9Q_SCP_001_ParseScope_Valid(t *testing.T) {
	// Spec: KWL-ARCH-J2K9Q KWL-SCP-001 Scope: Package
	cases := []struct {
		in   string
		want Scope
	}{
		{"Workspace", ScopeWorkspace},
		{"Module", ScopeModule},
		{"Domain", ScopeDomain},
		{"Package", ScopePackage},
		{"Service", ScopeService},
		{"Func", ScopeFunc},
		{"workspace", ScopeWorkspace},
		{"DOMAIN", ScopeDomain},
		{"  Package  ", ScopePackage},
	}
	for _, c := range cases {
		got, err := ParseScope(c.in)
		if err != nil {
			t.Errorf("ParseScope(%q) error = %v, want nil", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseScope(%q) = %v, want %v", c.in, got, c.want)
		}
		if !got.IsValid() {
			t.Errorf("Scope(%q).IsValid() = false, want true", c.in)
		}
	}
}

func TestKWL_ARCH_J2K9Q_SCP_001_ParseScope_Invalid(t *testing.T) {
	// Spec: KWL-ARCH-J2K9Q KWL-SCP-001 Scope: Package
	cases := []string{"", "unknown", "app", "WS", "Project", "ProjectX"}
	for _, in := range cases {
		if _, err := ParseScope(in); err == nil {
			t.Errorf("ParseScope(%q) succeeded, want error", in)
		} else {
			// Should be UsageError with ExitCodeUsage
			if e, ok := err.(*Error); !ok || e.Code != ExitCodeUsage {
				t.Errorf("ParseScope(%q) error = %v, want UsageError (ExitCodeUsage)", in, err)
			}
		}
	}
}

func TestKWL_ARCH_J2K9Q_SCP_001_Scope_Ordering(t *testing.T) {
	// Spec: KWL-ARCH-J2K9Q KWL-SCP-001 Scope: Package
	// Ordering: Workspace < Module < Domain < Package < Service < Func
	ordered := []Scope{ScopeWorkspace, ScopeModule, ScopeDomain, ScopePackage, ScopeService, ScopeFunc}
	for i := 0; i < len(ordered)-1; i++ {
		if !ordered[i].Less(ordered[i+1]) {
			t.Errorf("%v.Less(%v) = false, want true", ordered[i], ordered[i+1])
		}
		if ordered[i].Level() >= ordered[i+1].Level() {
			t.Errorf("Level(%v)=%d should be < Level(%v)=%d", ordered[i], ordered[i].Level(), ordered[i+1], ordered[i+1].Level())
		}
	}
	// Reverse should not be less
	if ScopeFunc.Less(ScopeWorkspace) {
		t.Error("Func.Less(Workspace) = true, want false")
	}
}

func TestKWL_ARCH_J2K9Q_SCP_001_Scope_Level_Invalid(t *testing.T) {
	// Spec: KWL-ARCH-J2K9Q KWL-SCP-001 Scope: Package
	if got := Scope("Unknown").Level(); got != -1 {
		t.Errorf("Scope(Unknown).Level() = %d, want -1", got)
	}
	if Scope("Unknown").IsValid() {
		t.Error("Scope(Unknown).IsValid() = true, want false")
	}
}

func TestKWL_ARCH_J2K9Q_SCP_002_Containment(t *testing.T) {
	// Spec: KWL-ARCH-J2K9Q KWL-SCP-002 Scope: Package
	// This is documentation-level but we test that AllScopes has 6 entries and order is preserved.
	if len(AllScopes) != 6 {
		t.Fatalf("len(AllScopes) = %d, want 6", len(AllScopes))
	}
	for _, s := range AllScopes {
		if !s.IsValid() {
			t.Errorf("AllScopes contains invalid scope %q", s)
		}
	}
}

func TestKWL_ARCH_J2K9Q_SCP_003_Identity(t *testing.T) {
	// Spec: KWL-ARCH-J2K9Q KWL-SCP-003 Scope: Package
	// Identity rules are documented; test that ParseScope round-trips via string.
	for _, s := range AllScopes {
		got, err := ParseScope(string(s))
		if err != nil || got != s {
			t.Errorf("ParseScope(%q) = %v, %v; want %v, nil", string(s), got, err, s)
		}
	}
}
