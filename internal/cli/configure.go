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
)

// ConfigureCommand returns a command for configuring the CLI with an API token.
func ConfigureCommand() *cli.Command {
	return &cli.Command{
		Name:  "configure",
		Usage: "Configure the CLI with your API token",
		Action: func(_ context.Context, _ *cli.Command) error {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter your Capacities API Token: ")
			token, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read token: %w", err)
			}
			token = strings.TrimSpace(token)

			fmt.Print("Enter your Default Space ID (optional): ")
			spaceID, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read space ID: %w", err)
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
