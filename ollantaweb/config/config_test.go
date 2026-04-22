package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesTomlConfigFile(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "ollanta.toml")
	t.Setenv("OLLANTA_CONFIG_FILE", configPath)
	if err := os.WriteFile(configPath, []byte(`
[server]
addr = ":9090"
database_url = "postgres://ollanta:ollanta_dev@localhost:5432/ollanta?sslmode=disable"
search_backend = "postgres"
log_level = "debug"
scanner_token = "scanner-token"
jwt_expiry = "30m"
refresh_expiry = "72h"
oauth_redirect_base = "http://localhost:9090"

[zincsearch]
url = "http://localhost:4081"
user = "zinc-user"
password = "zinc-pass"

[oauth.github]
client_id = "github-id"
client_secret = "github-secret"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Addr != ":9090" || cfg.DatabaseURL == "" || cfg.ZincSearchURL != "http://localhost:4081" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.SearchBackend != "postgres" || cfg.LogLevel != "debug" || cfg.ScannerToken != "scanner-token" {
		t.Fatalf("unexpected server values: %+v", cfg)
	}
	if cfg.GitHubClientID != "github-id" || cfg.GitHubClientSecret != "github-secret" {
		t.Fatalf("unexpected github oauth config: %+v", cfg)
	}
	if cfg.JWTExpiry.Minutes() != 30 || cfg.RefreshExpiry.Hours() != 72 {
		t.Fatalf("unexpected durations: jwt=%s refresh=%s", cfg.JWTExpiry, cfg.RefreshExpiry)
	}
}

func TestLoadEnvironmentOverridesToml(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "ollanta.toml")
	t.Setenv("OLLANTA_CONFIG_FILE", configPath)
	if err := os.WriteFile(configPath, []byte(`
[server]
database_url = "postgres://from-file"
log_level = "info"
`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("OLLANTA_DATABASE_URL", "postgres://from-env")
	t.Setenv("OLLANTA_LOG_LEVEL", "warn")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabaseURL != "postgres://from-env" || cfg.LogLevel != "warn" {
		t.Fatalf("expected environment override, got %+v", cfg)
	}
}

func TestLoadSupportsStructuredSharedSections(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "ollanta.toml")
	t.Setenv("OLLANTA_CONFIG_FILE", configPath)
	if err := os.WriteFile(configPath, []byte(`
[server]
host = "0.0.0.0"
port = 8181
log_level = "debug"
scanner_token = "scanner-token"

[database]
host = "db.internal"
port = 5544
name = "ollanta"
user = "ollanta"
password = "secret"
sslmode = "require"

[search]
backend = "zincsearch"
scheme = "http"
host = "search.internal"
port = 4180
user = "search-user"
password = "search-pass"
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Addr != "0.0.0.0:8181" {
		t.Fatalf("unexpected addr: %q", cfg.Addr)
	}
	wantDB := "postgres://ollanta:secret@db.internal:5544/ollanta?sslmode=require"
	if cfg.DatabaseURL != wantDB {
		t.Fatalf("unexpected database url: %q", cfg.DatabaseURL)
	}
	if cfg.ZincSearchURL != "http://search.internal:4180" {
		t.Fatalf("unexpected search url: %q", cfg.ZincSearchURL)
	}
	if cfg.ZincSearchUser != "search-user" || cfg.ZincSearchPassword != "search-pass" {
		t.Fatalf("unexpected search credentials: %+v", cfg)
	}
}
