// Package main is the entry point for the Capacities CLI application.
package main

import (
	"context"
	"fmt"
	"os"

	internalcli "github.com/klauern/capacities-cli/internal/cli"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "capacities",
		Usage: "CLI for Capacities.io",
		Commands: []*cli.Command{
			internalcli.DailyCommand(),
			internalcli.SpacesCommand(),
			internalcli.SpaceInfoCommand(),
			internalcli.SearchCommand(),
			internalcli.SaveWebLinkCommand(),
			internalcli.ConfigureCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
