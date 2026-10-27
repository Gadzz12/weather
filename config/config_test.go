package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()
	if cfg.AppPort == "" {
		t.Fatalf("expected AppPort to have a default value or load from environment, got empty string")
	}
	if cfg.AppPort != "8080" {
		t.Logf("AppPort is set to %s", cfg.AppPort)
	}
}
