package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/urfave/cli/v3"
)

// SpacesCommand returns a command for listing all spaces.
func SpacesCommand() *cli.Command {
	return &cli.Command{
		Name:  "spaces",
		Usage: "List all spaces",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format: table or json",
				Value: formatTable,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			format := cmd.String("format")
			if err := validateFormat(format); err != nil {
				return err
			}

			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			client := api.NewClient(cfg.Token)
			spaces, err := client.GetSpaces(ctx)
			if err != nil {
				return fmt.Errorf("failed to get spaces: %w", err)
			}

			return printSpaces(os.Stdout, spaces, format)
		},
	}
}
