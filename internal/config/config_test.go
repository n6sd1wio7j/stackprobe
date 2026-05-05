package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stackprobe/stackprobe/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "stackprobe-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempConfig(t, `{
		"listen_addr": ":9090",
		"interval": 10000000000,
		"endpoints": [
			{"name": "api", "url": "http://localhost/health", "timeout": 3000000000}
		]
	}`)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ListenAddr != ":9090" {
		t.Errorf("listen_addr: got %q, want %q", cfg.ListenAddr, ":9090")
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("interval: got %v, want %v", cfg.Interval, 10*time.Second)
	}
	if len(cfg.Endpoints) != 1 {
		t.Fatalf("endpoints: got %d, want 1", len(cfg.Endpoints))
	}
	if cfg.Endpoints[0].Name != "api" {
		t.Errorf("endpoint name: got %q, want %q", cfg.Endpoints[0].Name, "api")
	}
}

func TestLoad_DefaultTimeout(t *testing.T) {
	path := writeTempConfig(t, `{
		"endpoints": [{"name": "svc", "url": "http://svc/health"}]
	}`)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Endpoints[0].Timeout != 5*time.Second {
		t.Errorf("default timeout: got %v, want 5s", cfg.Endpoints[0].Timeout)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeTempConfig(t, `{invalid json}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoad_MissingEndpointName(t *testing.T) {
	path := writeTempConfig(t, `{
		"endpoints": [{"url": "http://svc/health"}]
	}`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing endpoint name, got nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.ListenAddr != ":8080" {
		t.Errorf("default listen_addr: got %q, want :8080", cfg.ListenAddr)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("default interval: got %v, want 30s", cfg.Interval)
	}
}
