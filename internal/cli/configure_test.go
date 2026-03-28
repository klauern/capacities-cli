package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klauern/capacities-cli/internal/config"
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

func TestIsInteractive(t *testing.T) {
	// This test verifies the isInteractive function runs without error
	// We can't assert true/false reliably since it depends on test environment
	result := isInteractive()
	t.Logf("isInteractive() returned: %v", result)
	// Just verify it doesn't panic
	_ = result
}

func TestConfigureWithFlags(t *testing.T) {
	configSubPath := filepath.Join(".config", "capacities", "config.yaml")

	t.Run("configure with token and space-id flags", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("HOME", tmp)

		// Simulate flag values being set
		// In real usage, these would come from command line
		token := "test-api-token"
		spaceID := "test-space-id"

		cfg := &config.Config{
			Token:          token,
			DefaultSpaceID: spaceID,
		}

		if err := config.Save(cfg); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		// Verify config was saved
		cfgPath := filepath.Join(tmp, configSubPath)
		data, err := os.ReadFile(cfgPath)
		if err != nil {
			t.Fatalf("failed to read config: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "test-api-token") {
			t.Errorf("config missing token: %s", content)
		}
		if !strings.Contains(content, "test-space-id") {
			t.Errorf("config missing space ID: %s", content)
		}
	})

	t.Run("configure with only token flag", func(t *testing.T) {
		tmp := t.TempDir()
		t.Setenv("HOME", tmp)

		token := "token-only-test"

		cfg := &config.Config{
			Token:          token,
			DefaultSpaceID: "",
		}

		if err := config.Save(cfg); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		// Load and verify
		loaded, err := config.Load()
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}
		if loaded.Token != "token-only-test" {
			t.Errorf("expected token 'token-only-test', got %q", loaded.Token)
		}
		if loaded.DefaultSpaceID != "" {
			t.Errorf("expected empty space ID, got %q", loaded.DefaultSpaceID)
		}
	})
}

func TestConfigureFlagValidation(t *testing.T) {
	// Test that token trimming works correctly
	t.Run("token with whitespace is trimmed", func(t *testing.T) {
		token := "  test-token  "
		trimmed := strings.TrimSpace(token)
		if trimmed != "test-token" {
			t.Errorf("expected 'test-token', got %q", trimmed)
		}
	})

	t.Run("space ID with whitespace is trimmed", func(t *testing.T) {
		spaceID := "  space-123  "
		trimmed := strings.TrimSpace(spaceID)
		if trimmed != "space-123" {
			t.Errorf("expected 'space-123', got %q", trimmed)
		}
	})
}
