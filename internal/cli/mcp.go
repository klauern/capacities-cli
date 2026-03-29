package cli

import (
	"context"
	"fmt"

	mcpserver "github.com/klauern/capacities-cli/internal/mcp"
	"github.com/urfave/cli/v3"
)

// MCPCommand returns commands for enabling and running the Capacities MCP server.
func MCPCommand() *cli.Command {
	return &cli.Command{
		Name:  "mcp",
		Usage: "Manage the Capacities MCP server",
		Commands: []*cli.Command{
			mcpServeCommand(),
		},
	}
}

func mcpServeCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run the Capacities MCP server over stdio",
		Action: func(ctx context.Context, _ *cli.Command) error {
			if err := mcpserver.Run(ctx); err != nil {
				return fmt.Errorf("failed to run MCP server: %w", err)
			}
			return nil
		},
	}
}
