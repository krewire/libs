package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCleanRejectsEscapingPath(t *testing.T) {
	root := t.TempDir()
	if _, err := Clean(root, "../outside.txt"); !errors.Is(err, os.ErrPermission) {
		t.Fatalf("Clean error = %v, want permission error", err)
	}
}

func TestCleanUsesRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root")
	got, err := Clean(root, "nested/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "nested", "file.txt")
	if got != want {
		t.Fatalf("Clean = %q, want %q", got, want)
	}
}
