package report

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCountLines_GoSource(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	src := `package main

// This is a comment
import "fmt"

/*
  block
*/
func main() {
	fmt.Println("hi")
}
`
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	total, ncloc, comments := countLines(path)
	if total != 11 {
		t.Errorf("total: got %d, want 11", total)
	}
	if comments != 4 {
		t.Errorf("comments: got %d, want 4", comments)
	}
	if ncloc != 5 {
		t.Errorf("ncloc: got %d, want 5", ncloc)
	}
}

func TestCountLines_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.go")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	total, ncloc, comments := countLines(path)
	if total != 0 || ncloc != 0 || comments != 0 {
		t.Errorf("empty: got (%d,%d,%d), want (0,0,0)", total, ncloc, comments)
	}
}

func TestCountLines_Unreadable(t *testing.T) {
	total, ncloc, comments := countLines("/nonexistent/file.go")
	if total != 0 || ncloc != 0 || comments != 0 {
		t.Errorf("unreadable: got (%d,%d,%d), want (0,0,0)", total, ncloc, comments)
	}
}

func TestCountLines_HashComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script.py")
	src := "# comment\ncode\n"
	if err := os.WriteFile(path, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	total, ncloc, comments := countLines(path)
	if total != 2 {
		t.Errorf("total: got %d, want 2", total)
	}
	if comments != 1 {
		t.Errorf("comments: got %d, want 1", comments)
	}
	if ncloc != 1 {
		t.Errorf("ncloc: got %d, want 1", ncloc)
	}
}
