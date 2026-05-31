// Package main is the entry point for the Capacities CLI application.
package main

import (
	"context"
	"fmt"
	"os"

	internalcli "github.com/klauern/capacities-cli/internal/cli"
	"github.com/urfave/cli/v3"
)

// Build-time variables injected via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// buildVersion returns a formatted version string including build metadata.
func buildVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, by: %s)", version, commit, date, builtBy)
}

func main() {
	cmd := &cli.Command{
		Name:    "capacities",
		Usage:   "CLI for Capacities.io",
		Version: buildVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Usage: "Capacities API token (overrides config; also via CAPACITIES_TOKEN)",
			},
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "Default space ID (overrides config; also via CAPACITIES_DEFAULT_SPACE_ID)",
			},
		},
		Commands: []*cli.Command{
			internalcli.DailyCommand(),
			internalcli.SpacesCommand(),
			internalcli.SpaceInfoCommand(),
			internalcli.SearchCommand(),
			internalcli.SaveWebLinkCommand(),
			internalcli.ConfigureCommand(),
			internalcli.MCPCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
