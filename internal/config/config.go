package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string        `json:"addr"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	LogLevel        string        `json:"log_level"`
	LogFormat       string        `json:"log_format"`
}

func Default() Config {
	return Config{
		Addr:            ":8080",
		ShutdownTimeout: 15 * time.Second,
		LogLevel:        "info",
		LogFormat:       "json",
	}
}

// Load returns config built from defaults, optionally overlaid by the JSON
// file at GIT_FORGE_CONFIG, and finally overlaid by environment variables.
func Load() (Config, error) {
	cfg := Default()

	if path := os.Getenv("GIT_FORGE_CONFIG"); path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return cfg, fmt.Errorf("read config %q: %w", path, err)
		}
		if err := json.Unmarshal(b, &cfg); err != nil {
			return cfg, fmt.Errorf("parse config %q: %w", path, err)
		}
	}

	if v := os.Getenv("GIT_FORGE_ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("GIT_FORGE_SHUTDOWN_TIMEOUT"); v != "" {
		secs, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("GIT_FORGE_SHUTDOWN_TIMEOUT: %w", err)
		}
		cfg.ShutdownTimeout = time.Duration(secs) * time.Second
	}
	if v := os.Getenv("GIT_FORGE_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("GIT_FORGE_LOG_FORMAT"); v != "" {
		cfg.LogFormat = v
	}

	return cfg, nil
}
