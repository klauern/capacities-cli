package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/klauern/capacities-cli/internal/config"
	"github.com/urfave/cli/v3"
)

func SearchCommand() *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "Search for content",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "ID of the space to search in (can be specified multiple times, comma separated)",
			},
			&cli.StringFlag{
				Name:  "mode",
				Usage: "Search mode (title or fullText)",
				Value: "title",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			term := strings.Join(cmd.Args().Slice(), " ")
			if term == "" {
				return fmt.Errorf("search term is required")
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if cfg.Token == "" {
				return fmt.Errorf("API token not found in config. Please configure it first")
			}

			spaceIDsStr := cmd.String("space-id")
			var spaceIDs []string
			if spaceIDsStr != "" {
				spaceIDs = strings.Split(spaceIDsStr, ",")
			} else if cfg.DefaultSpaceID != "" {
				spaceIDs = []string{cfg.DefaultSpaceID}
			}

			if len(spaceIDs) == 0 {
				return fmt.Errorf("at least one space ID is required (via flag or config)")
			}

			client := api.NewClient(cfg.Token)
			req := api.SearchRequest{
				SearchTerm: term,
				SpaceIDs:   spaceIDs,
				Mode:       cmd.String("mode"),
			}

			results, err := client.Search(ctx, req)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tTITLE\tSCORE")
			for _, r := range results {
				score := 0.0
				if len(r.Highlights) > 0 {
					score = r.Highlights[0].Score
				}
				fmt.Fprintf(w, "%s\t%s\t%.2f\n", r.ID, r.Title, score)
			}
			w.Flush()

			return nil
		},
	}
}
