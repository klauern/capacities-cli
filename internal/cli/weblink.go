package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/urfave/cli/v3"
)

// SaveWebLinkCommand returns a command for saving web links to Capacities.
func SaveWebLinkCommand() *cli.Command {
	return &cli.Command{
		Name:  "save-weblink",
		Usage: "Save a weblink to a space",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "title",
				Usage: "Overwrite title",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Overwrite description",
			},
			&cli.StringSliceFlag{
				Name:  "tags",
				Usage: "Tags to add",
			},
			&cli.StringFlag{
				Name:  "md-text",
				Usage: "Markdown text to add to the notes section",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			url := strings.Join(cmd.Args().Slice(), " ")
			if url == "" {
				return fmt.Errorf("URL is required")
			}

			auth, spaceID, err := RequireSpaceID(cmd)
			if err != nil {
				return err
			}

			client := api.NewClient(auth.Token)
			req := api.SaveWebLinkRequest{
				SpaceID:              spaceID,
				URL:                  url,
				TitleOverwrite:       cmd.String("title"),
				DescriptionOverwrite: cmd.String("description"),
				Tags:                 cmd.StringSlice("tags"),
				MDText:               cmd.String("md-text"),
			}

			resp, err := client.SaveWebLink(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to save weblink: %w", err)
			}

			fmt.Printf("Successfully saved weblink. ID: %s\n", resp.ID)
			return nil
		},
	}
}
