package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configSubPath := filepath.Join(".config", "capacities", "config.yaml")

	writeConfig := func(t *testing.T, content string) {
		t.Helper()
		tmp := t.TempDir()
		t.Setenv("HOME", tmp)
		cfgPath := filepath.Join(tmp, configSubPath)
		if err := os.MkdirAll(filepath.Dir(cfgPath), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(cfgPath, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("load error from malformed config", func(t *testing.T) {
		writeConfig(t, "token: [unclosed\n") // invalid YAML
		cfg, err := loadConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if cfg != nil {
			t.Fatalf("expected nil cfg, got %+v", cfg)
		}
	})

	t.Run("empty token returns error", func(t *testing.T) {
		writeConfig(t, "token: \"\"\n")
		cfg, err := loadConfig()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if cfg != nil {
			t.Fatalf("expected nil cfg, got %+v", cfg)
		}
		if !strings.Contains(err.Error(), "API token not found") {
			t.Fatalf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("valid token returns cfg", func(t *testing.T) {
		writeConfig(t, "token: my-token\ndefault_space_id: space-1\n")
		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if cfg == nil {
			t.Fatal("expected non-nil cfg")
		}
		if cfg.Token != "my-token" {
			t.Fatalf("expected token %q, got %q", "my-token", cfg.Token)
		}
		if cfg.DefaultSpaceID != "space-1" {
			t.Fatalf("expected default_space_id %q, got %q", "space-1", cfg.DefaultSpaceID)
		}
	})
}
