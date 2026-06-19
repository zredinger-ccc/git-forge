package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Addr != ":8080" {
		t.Errorf("default Addr = %q, want %q", c.Addr, ":8080")
	}
	if c.ShutdownTimeout != 15*time.Second {
		t.Errorf("default ShutdownTimeout = %v, want %v", c.ShutdownTimeout, 15*time.Second)
	}
}

func TestLoad_DefaultsOnly(t *testing.T) {
	t.Setenv("GIT_FORGE_CONFIG", "")
	t.Setenv("GIT_FORGE_ADDR", "")
	t.Setenv("GIT_FORGE_LOG_LEVEL", "")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.Addr != ":8080" {
		t.Errorf("Addr = %q, want %q", c.Addr, ":8080")
	}
}

func TestLoad_FileThenEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"addr":":9000","log_level":"debug"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_FORGE_CONFIG", path)
	t.Setenv("GIT_FORGE_ADDR", ":9999")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.Addr != ":9999" {
		t.Errorf("Addr = %q, want %q (env should win over file)", c.Addr, ":9999")
	}
	if c.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q (from file)", c.LogLevel, "debug")
	}
}

func TestLoad_BadShutdownTimeout(t *testing.T) {
	t.Setenv("GIT_FORGE_SHUTDOWN_TIMEOUT", "not-a-number")
	if _, err := Load(); err == nil {
		t.Fatal("expected error, got nil")
	}
}
