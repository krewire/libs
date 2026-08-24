// Tests for KWL-P8W2N
package core

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// Spec: KWL-P8W2N KWL-ERRV-008 Scope: Package
func TestKWL_ERRV_008_WithAttrs_ExtractsThroughWrapsKeepsIdentity(t *testing.T) {
	base := errors.New("disk full")
	wrapped := WithAttrs(WithStack(WithAttrs(base, Attr{"path", "/tmp/a.yaml"})), Attr{"step", "build"})

	got := AttrsOf(wrapped)
	if len(got) != 2 {
		t.Fatalf("AttrsOf len = %d (%v), want 2 attrs outermost-first", len(got), got)
	}
	if got[0].Key != "step" || got[1].Key != "path" {
		t.Errorf("AttrsOf order = %v/%v, want step then path", got[0], got[1])
	}

	var target *Error
	if !errors.As(wrapped, &target) && !errors.Is(wrapped, base) {
		t.Error("errors.Is/As broken through attr wrap")
	}
	if !errors.Is(wrapped, base) {
		t.Error("errors.Is(wrapped, base) = false, want true")
	}
	if AttrsOf(nil) != nil {
		t.Error("AttrsOf(nil) should be nil")
	}
	if WithAttrs(nil, Attr{"k", "v"}) != nil {
		t.Error("WithAttrs(nil) should stay nil")
	}
}

// Spec: KWL-P8W2N KWL-ERRV-009 Scope: Package
func TestKWL_ERRV_009_HintOf_ReturnsNearestHintThroughChain(t *testing.T) {
	base := FailureError("config invalid")
	deep := WithHint(base, "check project.kind in krewire.yaml")
	wrapped := WithHint(deep, "run 'kiw init' to create one")

	if got := HintOf(wrapped); got != "run 'kiw init' to create one" {
		t.Errorf("HintOf(nearest) = %q", got)
	}
	if got := HintOf(deep); got != "check project.kind in krewire.yaml" {
		t.Errorf("HintOf(deep) = %q", got)
	}
	if got := HintOf(base); got != "" {
		t.Errorf("HintOf(unhinted) = %q, want empty", got)
	}
	if got := HintOf(nil); got != "" {
		t.Errorf("HintOf(nil) = %q, want empty", got)
	}
}

// Spec: KWL-P8W2N KWL-ERRV-010 Scope: Package
func TestKWL_ERRV_010_FormatTree_RendersChainOriginAndHint(t *testing.T) {
	inner := errors.New("open krewire.yaml: no such file or directory")
	err := WithStack(fmt.Errorf("cannot read config: %w",
		WithHint(
			WithAttrs(inner, Attr{Key: "file", Value: "krewire.yaml"}),
			"run 'kiw init' to create one",
		),
	))

	tree := FormatTree(err)
	for _, want := range []string{
		"Error: cannot read config: open krewire.yaml",
		"file=krewire.yaml",
		"caused by: open krewire.yaml: no such file or directory",
		"caused by:",
		"at ",
		".go:",
		"Hint: run 'kiw init' to create one",
	} {
		if !strings.Contains(tree, want) {
			t.Errorf("tree missing %q:\n%s", want, tree)
		}
	}

	// A hintless plain error still renders a single-line tree without footer.
	plain := FormatTree(errors.New("boom"))
	if !strings.HasPrefix(plain, "Error: boom\n") || strings.Contains(plain, "Hint:") {
		t.Errorf("plain tree malformed:\n%s", plain)
	}
	if FormatTree(nil) != "" {
		t.Error("FormatTree(nil) should be empty")
	}
}
