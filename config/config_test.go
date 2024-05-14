package config

import "testing"

func TestConfig(t *testing.T) {
	cfg := LoadConfig("config.yml")
	expected := 2
	if cfg.Concurrency != expected {
		t.Errorf("Config read incorrect: %d, want: %d.", cfg.Concurrency, expected)
	}
}
