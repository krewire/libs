// Tests for KWL-P8W2N
package core

import (
	"errors"
	"strings"
	"testing"
)

// Spec: KWL-P8W2N KWL-ERRV-001 S3 Scope: Unit
func TestKWL_ERRV_001_WithStack_PreservesIdentityAndUnwrap(t *testing.T) {
	target := UsageError("boom")
	wrapped := WithStack(WithStack(target))

	if !errors.Is(wrapped, target) {
		t.Error("errors.Is through two stacks = false, want true")
	}
	var ce interface{ ExitCode() ExitCode }
	if !errors.As(wrapped, &ce) || ce.ExitCode() != ExitCodeUsage {
		t.Error("errors.As through stack lost the usage code")
	}
	if WithStack(nil) != nil {
		t.Error("WithStack(nil) != nil")
	}
}

// Spec: KWL-P8W2N KWL-ERRV-002 Scope: Unit
func TestKWL_ERRV_002_StackOf_ExtractsFramesAndFormats(t *testing.T) {
	mk := func() error { return WithStack(UsageError("origin")) }
	err := mk()

	frames := StackOf(err)
	if len(frames) == 0 {
		t.Fatal("StackOf = 0 frames, want the capture chain")
	}
	foundSelf := false
	for _, f := range frames {
		if strings.Contains(f.Func, "TestKWL_ERRV_002") && strings.Contains(f.File, "stack_test.go") && f.Line > 0 {
			foundSelf = true
		}
	}
	if !foundSelf {
		t.Errorf("frames missing the creation site: %v", frames)
	}
	rendered := FormatStack(frames)
	if !strings.Contains(rendered, "at ") || !strings.Contains(rendered, "stack_test.go:") {
		t.Errorf("FormatStack output incomplete:\n%s", rendered)
	}

	if got := StackOf(UsageError("no stack")); got != nil {
		t.Errorf("StackOf(unstacked) = %v, want nil", got)
	}
}

// Spec: KWL-P8W2N KWL-ERRV-003 Scope: Unit
func TestKWL_ERRV_003_DoubleStack_KeepsBothTraces(t *testing.T) {
	inner := WithStack(UsageError("inner"))
	outer := WithStack(inner)

	if got := len(StackOf(inner)); got == 0 {
		t.Error("inner trace lost")
	}
	if got := len(StackOf(outer)); got < len(StackOf(inner)) {
		t.Errorf("outer trace (%d frames) shallower than inner (%d)", len(StackOf(outer)), len(StackOf(inner)))
	}
}
