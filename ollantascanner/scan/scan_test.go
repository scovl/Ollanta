package scan_test

import (
	"testing"

	"github.com/scovl/ollanta/ollantascanner/scan"
)

func TestParseFlags_Defaults(t *testing.T) {
	opts, err := scan.ParseFlags([]string{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.ProjectDir != "." {
		t.Errorf("ProjectDir: got %q, want %q", opts.ProjectDir, ".")
	}
	if opts.Format != "all" {
		t.Errorf("Format: got %q, want %q", opts.Format, "all")
	}
	if opts.Debug {
		t.Error("Debug should be false by default")
	}
	if len(opts.Sources) == 0 {
		t.Error("Sources should have a default")
	}
}

func TestParseFlags_ProjectDir(t *testing.T) {
	opts, err := scan.ParseFlags([]string{"-project-dir", "/tmp/myproject"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.ProjectDir != "/tmp/myproject" {
		t.Errorf("ProjectDir: got %q", opts.ProjectDir)
	}
}

func TestParseFlags_ProjectKey(t *testing.T) {
	opts, err := scan.ParseFlags([]string{"-project-key", "myapp"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.ProjectKey != "myapp" {
		t.Errorf("ProjectKey: got %q", opts.ProjectKey)
	}
}

func TestParseFlags_Sources(t *testing.T) {
	opts, err := scan.ParseFlags([]string{"-sources", "./cmd/...,./pkg/..."})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.Sources) != 2 {
		t.Errorf("Sources: got %v", opts.Sources)
	}
}

func TestParseFlags_Exclusions(t *testing.T) {
	opts, err := scan.ParseFlags([]string{"-exclusions", "*_test.go,vendor/**"})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.Exclusions) != 2 {
		t.Errorf("Exclusions: got %v", opts.Exclusions)
	}
}

func TestParseFlags_Debug(t *testing.T) {
	opts, err := scan.ParseFlags([]string{"-debug"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.Debug {
		t.Error("Debug should be true")
	}
}

func TestParseFlags_DefaultProjectKey(t *testing.T) {
	// When -project-key is not set, ProjectKey defaults to the base name of ProjectDir.
	opts, err := scan.ParseFlags([]string{"-project-dir", "/some/path/myrepo"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.ProjectKey != "myrepo" {
		t.Errorf("ProjectKey default: got %q, want %q", opts.ProjectKey, "myrepo")
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := scan.ParseFlags([]string{"-unknown-flag"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := scan.ParseFlags([]string{"-format", "xml"})
	if err == nil {
		t.Error("expected error for unsupported format 'xml'")
	}
}

func TestParseFlags_ValidFormats(t *testing.T) {
	for _, f := range []string{"summary", "json", "sarif", "all"} {
		cfg, err := scan.ParseFlags([]string{"-format", f})
		if err != nil {
			t.Errorf("format %q should be valid: %v", f, err)
		}
		if cfg.Format != f {
			t.Errorf("format: got %q, want %q", cfg.Format, f)
		}
	}
}
