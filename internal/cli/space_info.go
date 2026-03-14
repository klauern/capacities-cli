package cli

import (
	"context"
	"fmt"
	"os"

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
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format: table or json",
				Value: formatTable,
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

			return printStructures(os.Stdout, structures, cmd.String("format"))
		},
	}
}
