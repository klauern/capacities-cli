package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v3"
)

func writeTestConfig(t *testing.T, homeDir string, token string, defaultSpaceID string) {
	t.Helper()

	configPath := filepath.Join(homeDir, ".config", "capacities", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	data := []byte("token: " + token + "\n" + "default_space_id: " + defaultSpaceID + "\n")
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func runResolveAuth(t *testing.T, args []string) (ResolvedAuth, error) {
	t.Helper()

	var got ResolvedAuth
	cmd := &cli.Command{
		Name: "capacities",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
			&cli.StringFlag{Name: "space-id"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			var err error
			got, err = ResolveAuth(cmd)
			return err
		},
	}

	runArgs := append([]string{"capacities"}, args...)
	err := cmd.Run(context.Background(), runArgs)
	return got, err
}

func TestResolveAuth_PrefersFlagOverEnvOverConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeTestConfig(t, home, "cfg-token", "cfg-space")

	t.Setenv("CAPACITIES_TOKEN", "env-token")
	t.Setenv("CAPACITIES_DEFAULT_SPACE_ID", "env-space")

	got, err := runResolveAuth(t, []string{"--token", "flag-token", "--space-id", "flag-space"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if got.Token != "flag-token" {
		t.Fatalf("token: got %q, want %q", got.Token, "flag-token")
	}
	if got.DefaultSpaceID != "flag-space" {
		t.Fatalf("space: got %q, want %q", got.DefaultSpaceID, "flag-space")
	}
}

func TestResolveAuth_UsesEnvWhenNoFlag(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeTestConfig(t, home, "cfg-token", "cfg-space")

	t.Setenv("CAPACITIES_TOKEN", "env-token")
	t.Setenv("CAPACITIES_DEFAULT_SPACE_ID", "env-space")

	got, err := runResolveAuth(t, nil)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if got.Token != "env-token" {
		t.Fatalf("token: got %q, want %q", got.Token, "env-token")
	}
	if got.DefaultSpaceID != "env-space" {
		t.Fatalf("space: got %q, want %q", got.DefaultSpaceID, "env-space")
	}
}

func TestResolveAuth_UsesConfigWhenNoFlagOrEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeTestConfig(t, home, "cfg-token", "cfg-space")

	got, err := runResolveAuth(t, nil)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if got.Token != "cfg-token" {
		t.Fatalf("token: got %q, want %q", got.Token, "cfg-token")
	}
	if got.DefaultSpaceID != "cfg-space" {
		t.Fatalf("space: got %q, want %q", got.DefaultSpaceID, "cfg-space")
	}
}

func TestRequireToken_ErrWhenMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	var gotErr error
	cmd := &cli.Command{
		Name: "capacities",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token"},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			_, gotErr = RequireToken(cmd)
			return nil
		},
	}
	_ = cmd.Run(context.Background(), []string{"capacities"})
	if gotErr == nil {
		t.Fatalf("expected error")
	}
}
