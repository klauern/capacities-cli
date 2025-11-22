package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/klauern/capacities-cli/internal/config"
	"github.com/urfave/cli/v3"
)

// SpacesCommand returns a command for listing all spaces.
func SpacesCommand() *cli.Command {
	return &cli.Command{
		Name:  "spaces",
		Usage: "List all spaces",
		Action: func(ctx context.Context, _ *cli.Command) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if cfg.Token == "" {
				return fmt.Errorf("API token not found in config. Please configure it first")
			}

			client := api.NewClient(cfg.Token)
			spaces, err := client.GetSpaces(ctx)
			if err != nil {
				return fmt.Errorf("failed to get spaces: %w", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			_, _ = fmt.Fprintln(w, "ID\tTITLE\tICON")
			for _, space := range spaces {
				iconVal := space.Icon.Val
				if space.Icon.Type == "emoji" {
					iconVal = space.Icon.Val
				}
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", space.ID, space.Title, iconVal)
			}
			_ = w.Flush()

			return nil
		},
	}
}
