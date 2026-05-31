package mcpserver

import (
	"context"
	"regexp"
	"testing"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/klauern/capacities-cli/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestServerToolNamesAreOpenAICompatible(t *testing.T) {
	srv, err := NewServer(&config.Config{Token: "test-token"}, api.NewClient("test-token"))
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	serverSession, err := srv.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("Connect(server) error = %v", err)
	}
	defer serverSession.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("Connect(client) error = %v", err)
	}
	defer clientSession.Close()

	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}

	toolNamePattern := regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	for _, tool := range tools.Tools {
		if !toolNamePattern.MatchString(tool.Name) {
			t.Fatalf("tool name %q is not OpenAI-compatible", tool.Name)
		}
		if len(tool.Name) > 64 {
			t.Fatalf("tool name %q exceeds OpenAI's 64-character limit", tool.Name)
		}
	}
}
