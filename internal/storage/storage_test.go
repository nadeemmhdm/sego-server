package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveRejectsEscape(t *testing.T) {
	d := t.TempDir()
	m, e := New([]Root{{ID: "r", Name: "R", Path: d}})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = m.Resolve("r", "../secret", false); e == nil {
		t.Fatal("escape accepted")
	}
}
func TestStreamAndList(t *testing.T) {
	d := t.TempDir()
	m, _ := New([]Root{{ID: "r", Name: "R", Path: d}})
	if _, e := m.StreamWrite("r", "a/b.txt", strings.NewReader("hello"), 100); e != nil {
		t.Fatal(e)
	}
	xs, e := m.List("r", "a")
	if e != nil || len(xs) != 1 || xs[0].Name != "b.txt" {
		t.Fatalf("bad list: %#v %v", xs, e)
	}
}
func TestSymlinkRejected(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permission varies")
	}
	d := t.TempDir()
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "x"), []byte("x"), 0600)
	os.Symlink(outside, filepath.Join(d, "link"))
	m, _ := New([]Root{{ID: "r", Name: "R", Path: d}})
	if _, e := m.Resolve("r", "link/x", true); e == nil {
		t.Fatal("symlink accepted")
	}
}
