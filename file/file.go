// Package file provides file-oriented operations built on the fs boundary.
package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/krewire/libs/fs"
)

const defaultFileMode os.FileMode = 0o600

// Read returns the complete contents of name.
func Read(fsys fs.FileSystem, name string) ([]byte, error) {
	if fsys == nil {
		fsys = fs.Default
	}
	data, err := fsys.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("file: read %q: %w", name, err)
	}
	return data, nil
}

// Write writes data to name, creating or replacing the file.
func Write(fsys fs.FileSystem, name string, data []byte, perm os.FileMode) error {
	if fsys == nil {
		fsys = fs.Default
	}
	if perm == 0 {
		perm = defaultFileMode
	}
	if err := fsys.MkdirAll(filepath.Dir(name), 0o700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("file: create parent for %q: %w", name, err)
	}
	if err := fsys.WriteFile(name, data, perm); err != nil {
		return fmt.Errorf("file: write %q: %w", name, err)
	}
	return nil
}

// WriteAtomic writes data through a temporary file and renames it into place.
// The temporary file is created beside the destination so the rename remains
// on one filesystem.
func WriteAtomic(fsys fs.FileSystem, name string, data []byte, perm os.FileMode) error {
	if fsys == nil {
		fsys = fs.Default
	}
	if perm == 0 {
		perm = defaultFileMode
	}
	if err := fsys.MkdirAll(filepath.Dir(name), 0o700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("file: create parent for %q: %w", name, err)
	}
	tmp := fmt.Sprintf("%s.%d.tmp", name, time.Now().UnixNano())
	if err := fsys.WriteFile(tmp, data, perm); err != nil {
		return fmt.Errorf("file: write temporary %q: %w", tmp, err)
	}
	if err := fsys.Rename(tmp, name); err != nil {
		_ = fsys.Remove(tmp)
		return fmt.Errorf("file: replace %q: %w", name, err)
	}
	return nil
}

// Copy copies all bytes from src to dst.
func Copy(fsys fs.FileSystem, dst, src string, perm os.FileMode) error {
	if fsys == nil {
		fsys = fs.Default
	}
	in, err := fsys.Open(src)
	if err != nil {
		return fmt.Errorf("file: open %q: %w", src, err)
	}
	defer in.Close()
	data, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("file: read %q: %w", src, err)
	}
	return Write(fsys, dst, data, perm)
}
