package file

import (
	"testing"

	"github.com/krewire/libs/fs"
)

func TestWriteAtomicReplacesContents(t *testing.T) {
	path := t.TempDir() + "/nested/data.txt"
	if err := WriteAtomic(fs.Default, path, []byte("first"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteAtomic(fs.Default, path, []byte("second"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := Read(fs.Default, path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Fatalf("contents = %q, want second", got)
	}
}
