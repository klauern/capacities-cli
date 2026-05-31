package cli

import (
	"fmt"
	"os"

	"github.com/klauern/capacities-cli/internal/config"
	"github.com/urfave/cli/v3"
)

const (
	envToken          = "CAPACITIES_TOKEN"
	envDefaultSpaceID = "CAPACITIES_DEFAULT_SPACE_ID"
)

// ResolvedAuth contains the token and default space ID chosen for a command.
type ResolvedAuth struct {
	Token          string
	DefaultSpaceID string
}

// ResolveAuth resolves Capacities credentials from flags, environment, or config.
func ResolveAuth(cmd *cli.Command) (ResolvedAuth, error) {
	cfg, err := config.Load()
	if err != nil {
		return ResolvedAuth{}, fmt.Errorf("failed to load config: %w", err)
	}

	token := cmd.String("token")
	if token == "" {
		token = os.Getenv(envToken)
	}
	if token == "" {
		token = cfg.Token
	}

	defaultSpaceID := cmd.String("space-id")
	if defaultSpaceID == "" {
		defaultSpaceID = os.Getenv(envDefaultSpaceID)
	}
	if defaultSpaceID == "" {
		defaultSpaceID = cfg.DefaultSpaceID
	}

	return ResolvedAuth{
		Token:          token,
		DefaultSpaceID: defaultSpaceID,
	}, nil
}

// RequireToken resolves auth and returns an error when no API token is available.
func RequireToken(cmd *cli.Command) (ResolvedAuth, error) {
	auth, err := ResolveAuth(cmd)
	if err != nil {
		return ResolvedAuth{}, err
	}
	if auth.Token == "" {
		return ResolvedAuth{}, fmt.Errorf("API token is required (set --token, CAPACITIES_TOKEN, or run 'capacities configure')")
	}
	return auth, nil
}

// RequireSpaceID resolves auth and returns an error when no default space ID is available.
func RequireSpaceID(cmd *cli.Command) (ResolvedAuth, error) {
	auth, err := RequireToken(cmd)
	if err != nil {
		return ResolvedAuth{}, err
	}
	if auth.DefaultSpaceID == "" {
		return ResolvedAuth{}, fmt.Errorf("space ID is required (set --space-id, CAPACITIES_DEFAULT_SPACE_ID, or configure default_space_id)")
	}
	return auth, nil
}
