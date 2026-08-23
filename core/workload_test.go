package core

import "testing"

func TestParseKind(t *testing.T) {
	cases := []struct {
		in   string
		want Kind
		ok   bool
	}{
		{"app", KindApp, true},
		{"cli", KindCLI, true},
		{"worker", KindWorker, true},
		{"service", KindService, true},
		{"infra", KindInfra, true},
		{"kernel", KindKernel, true},
		{"unknown", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, err := ParseKind(c.in)
		if c.ok && err != nil {
			t.Errorf("ParseKind(%q) error = %v, want nil", c.in, err)
		}
		if !c.ok && err == nil {
			t.Errorf("ParseKind(%q) = %v, want error", c.in, got)
		}
		if c.ok && got != c.want {
			t.Errorf("ParseKind(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestWorkloadFor(t *testing.T) {
	w, ok := WorkloadFor(KindCLI)
	if !ok {
		t.Fatal("WorkloadFor(cli) not found")
	}
	if w.Package != "framework/tui" {
		t.Errorf("Package = %q, want %q", w.Package, "framework/tui")
	}
	_, ok = WorkloadFor(Kind("unknown"))
	if ok {
		t.Error("WorkloadFor(unknown) should not be found")
	}
}

func TestKindIsValid(t *testing.T) {
	if !KindApp.IsValid() {
		t.Error("KindApp should be valid")
	}
	if Kind("nope").IsValid() {
		t.Error("unknown kind should be invalid")
	}
}
