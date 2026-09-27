package term

import "testing"

func TestPaint(t *testing.T) {
	got := Paint("ok", ColorGreen, []Style{StyleBold})
	want := "\x1b[1;32mok\x1b[0m"
	if got != want {
		t.Errorf("Paint() = %q, want %q", got, want)
	}
}

func TestTerminalPaintWithoutColor(t *testing.T) {
	term := &Terminal{ColorSupported: false}
	if got := term.Paint("ok", ColorGreen, []Style{StyleBold}); got != "ok" {
		t.Errorf("Paint() = %q, want %q", got, "ok")
	}
}

func TestNewTerminal(t *testing.T) {
	// Assert only that construction never panics; the detected value depends
	// on the environment.
	_ = NewTerminal()
}

func TestNewTerminalFor_NilAndRegularFile(t *testing.T) {
	if IsTerminal(nil) {
		t.Error("nil file should not be terminal")
	}
	term := NewTerminalFor(nil)
	if term.ColorSupported {
		t.Error("nil file should not have ColorSupported without force env")
	}
}

func TestNewTerminalFor_ForceColor(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")
	term := NewTerminalFor(nil)
	if !term.ColorSupported {
		t.Error("CLICOLOR_FORCE=1 should enable color even for non-tty")
	}
}
