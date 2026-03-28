// Package cli provides command-line interface commands for the Capacities CLI.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/klauern/capacities-cli/internal/config"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// loadConfig loads the CLI configuration and validates that an API token is present.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("API token not found in config. Please configure it first")
	}
	return cfg, nil
}

// readTokenInteractive reads the API token from stdin without echoing (for interactive use).
func readTokenInteractive() (string, error) {
	fmt.Print("Enter your Capacities API Token: ")
	tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}
	fmt.Println() // Newline after password input
	return string(tokenBytes), nil
}

// readSpaceIDInteractive reads the optional space ID from stdin (echoed).
func readSpaceIDInteractive() (string, error) {
	fmt.Print("Enter your Default Space ID (optional): ")
	reader := bufio.NewReader(os.Stdin)
	spaceID, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read space ID: %w", err)
	}
	return strings.TrimSpace(spaceID), nil
}

// isInteractive returns true if stdin is a terminal.
func isInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// ConfigureCommand returns a command for configuring the CLI with an API token.
func ConfigureCommand() *cli.Command {
	return &cli.Command{
		Name:  "configure",
		Usage: "Configure the CLI with your API token",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Usage: "API token (if not provided, will prompt interactively)",
			},
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "Default space ID (optional, will prompt interactively if not provided)",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			var token, spaceID string
			var err error

			// Handle token input
			token = cmd.String("token")
			if token == "" {
				if isInteractive() {
					token, err = readTokenInteractive()
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("--token is required when running non-interactively")
				}
			}
			token = strings.TrimSpace(token)

			// Handle space ID input
			spaceID = cmd.String("space-id")
			if spaceID == "" && isInteractive() {
				spaceID, err = readSpaceIDInteractive()
				if err != nil {
					return err
				}
			}
			spaceID = strings.TrimSpace(spaceID)

			cfg := &config.Config{
				Token:          token,
				DefaultSpaceID: spaceID,
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Println("Configuration saved successfully!")
			return nil
		},
	}
}
