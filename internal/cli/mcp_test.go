package cli

import "testing"

func TestMCPCommandRegistration(t *testing.T) {
	cmd := MCPCommand()
	if cmd.Name != "mcp" {
		t.Fatalf("expected command name %q, got %q", "mcp", cmd.Name)
	}
	if len(cmd.Commands) != 1 {
		t.Fatalf("expected 1 subcommand, got %d", len(cmd.Commands))
	}
	if cmd.Commands[0].Name != "serve" {
		t.Fatalf("expected subcommand %q, got %q", "serve", cmd.Commands[0].Name)
	}
}
