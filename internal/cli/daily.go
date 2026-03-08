package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/urfave/cli/v3"
)

// DailyCommand returns a command for interacting with daily notes.
func DailyCommand() *cli.Command {
	return &cli.Command{
		Name:  "daily",
		Usage: "Interact with daily notes",
		Commands: []*cli.Command{
			{
				Name:  "save",
				Usage: "Save text to today's daily note",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "space-id",
						Usage: "ID of the space to save to",
					},
					&cli.BoolFlag{
						Name:  "no-timestamp",
						Usage: "Do not add a timestamp to the note",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					text := strings.Join(cmd.Args().Slice(), " ")
					if text == "" {
						return fmt.Errorf("text to save is required")
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
					opts := api.SaveOptions{
						NoTimeStamp: cmd.Bool("no-timestamp"),
					}

					if err := client.SaveToDailyNote(ctx, spaceID, text, opts); err != nil {
						return fmt.Errorf("failed to save to daily note: %w", err)
					}

					fmt.Println("Successfully saved to daily note")
					return nil
				},
			},
		},
	}
}
