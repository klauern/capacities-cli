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

			auth, err := RequireSpaceID(cmd)
			if err != nil {
				return err
			}

			client := api.NewClient(auth.Token)
			structures, err := client.GetSpaceInfo(ctx, auth.DefaultSpaceID)
			if err != nil {
				return fmt.Errorf("failed to get space info: %w", err)
			}

			return printStructures(os.Stdout, structures, format)
		},
	}
}
