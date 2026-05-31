package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/urfave/cli/v3"
)

func resolveDailySaveText(args []string, stdin io.Reader, stdinIsTTY bool, filePath string) (string, error) {
	if filePath != "" && len(args) > 0 {
		return "", fmt.Errorf("cannot use --file with positional text arguments")
	}

	if filePath != "" {
		var content []byte
		var err error
		if filePath == "-" {
			content, err = io.ReadAll(stdin)
		} else {
			// #nosec G304 -- --file intentionally reads the user-specified note content path.
			content, err = os.ReadFile(filePath)
		}
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		text := strings.TrimRight(string(content), "\r\n")
		if text == "" {
			return "", fmt.Errorf("text to save is required")
		}
		return text, nil
	}

	if len(args) > 0 {
		text := strings.Join(args, " ")
		if text == "" {
			return "", fmt.Errorf("text to save is required")
		}
		return text, nil
	}

	if !stdinIsTTY {
		content, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}
		text := strings.TrimRight(string(content), "\r\n")
		if text == "" {
			return "", fmt.Errorf("text to save is required")
		}
		return text, nil
	}

	return "", fmt.Errorf("text to save is required (provide args, --file, or stdin)")
}

func isTerminalStdin() (bool, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}
	return (stat.Mode() & os.ModeCharDevice) != 0, nil
}

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
						Name:  "file",
						Usage: "Read text to save from a file (use '-' for stdin)",
					},
					&cli.BoolFlag{
						Name:  "no-timestamp",
						Usage: "Do not add a timestamp to the note",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					stdinIsTTY, err := isTerminalStdin()
					if err != nil {
						return fmt.Errorf("failed to inspect stdin: %w", err)
					}

					text, err := resolveDailySaveText(cmd.Args().Slice(), os.Stdin, stdinIsTTY, cmd.String("file"))
					if err != nil {
						return err
					}

					auth, err := RequireSpaceID(cmd)
					if err != nil {
						return err
					}

					client := api.NewClient(auth.Token)
					opts := api.SaveOptions{
						NoTimeStamp: cmd.Bool("no-timestamp"),
					}

					if err := client.SaveToDailyNote(ctx, auth.DefaultSpaceID, text, opts); err != nil {
						return fmt.Errorf("failed to save to daily note: %w", err)
					}

					fmt.Println("Successfully saved to daily note")
					return nil
				},
			},
		},
	}
}
