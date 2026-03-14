package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/urfave/cli/v3"
)

// SearchCommand returns a command for looking up content by title in Capacities.
func SearchCommand() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Look up content by title",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "ID of the space to search in",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format: table or json",
				Value: formatTable,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			term := strings.Join(cmd.Args().Slice(), " ")
			if term == "" {
				return fmt.Errorf("search term is required")
			}

			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			spaceID := cmd.String("space-id")
			if spaceID == "" {
				spaceID = cfg.DefaultSpaceID
			}
			if spaceID == "" {
				return fmt.Errorf("space ID is required (either via flag or config)")
			}

			client := api.NewClient(cfg.Token)
			req := api.LookupRequest{
				SearchTerm: term,
				SpaceID:    spaceID,
			}

			results, err := client.Lookup(ctx, req)
			if err != nil {
				return fmt.Errorf("lookup failed: %w", err)
			}

			return printLookupResults(os.Stdout, results, cmd.String("format"))
		},
	}
}
