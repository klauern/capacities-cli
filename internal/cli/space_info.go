package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/urfave/cli/v3"
)

// SpaceInfoCommand returns a command for retrieving space information.
func SpaceInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "space-info",
		Usage: "Get info about a space (structures and collections)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "space-id",
				Usage: "ID of the space",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
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
			structures, err := client.GetSpaceInfo(ctx, spaceID)
			if err != nil {
				return fmt.Errorf("failed to get space info: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			_, _ = fmt.Fprintln(w, "STRUCTURE ID\tTITLE\tCOLLECTIONS")
			for _, s := range structures {
				var collections []string
				for _, c := range s.Collections {
					collections = append(collections, c.Title)
				}
				cols := strings.Join(collections, ", ")
				if cols == "" {
					cols = "-"
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", s.ID, s.Title, cols)
			}
			_ = w.Flush()

			return nil
		},
	}
}
