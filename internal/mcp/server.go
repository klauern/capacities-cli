// Package mcpserver provides the Capacities MCP server implementation.
package mcpserver

import (
	"context"
	"fmt"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/klauern/capacities-cli/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const serverVersion = "dev"

// Server exposes Capacities API operations as MCP tools.
type Server struct {
	cfg    *config.Config
	client *api.Client
	server *mcp.Server
}

// NewServer constructs a Capacities MCP server from the provided config.
func NewServer(cfg *config.Config, client *api.Client) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("API token not found in config. Please configure it first")
	}
	if client == nil {
		client = api.NewClient(cfg.Token)
	}

	srv := &Server{
		cfg:    cfg,
		client: client,
		server: mcp.NewServer(&mcp.Implementation{Name: "capacities-cli", Version: serverVersion}, nil),
	}

	mcp.AddTool(srv.server, &mcp.Tool{
		Name:        "save_to_daily_note",
		Description: "Save text to today's daily note",
	}, srv.saveToDailyNote)
	mcp.AddTool(srv.server, &mcp.Tool{
		Name:        "search",
		Description: "Look up content by title in a space",
	}, srv.search)
	mcp.AddTool(srv.server, &mcp.Tool{
		Name:        "list_spaces",
		Description: "List all spaces",
	}, srv.listSpaces)
	mcp.AddTool(srv.server, &mcp.Tool{
		Name:        "space_info",
		Description: "Get detailed information about a space",
	}, srv.spaceInfo)
	mcp.AddTool(srv.server, &mcp.Tool{
		Name:        "save_weblink",
		Description: "Save a web link to a space",
	}, srv.saveWebLink)

	return srv, nil
}

// Run starts the MCP server over stdio.
func (s *Server) Run(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

// Connect connects the server to a transport, primarily for tests.
func (s *Server) Connect(ctx context.Context, transport mcp.Transport, opts *mcp.ServerSessionOptions) (*mcp.ServerSession, error) {
	return s.server.Connect(ctx, transport, opts)
}

// Run loads the default config and starts the MCP server over stdio.
func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	srv, err := NewServer(cfg, nil)
	if err != nil {
		return err
	}

	return srv.Run(ctx)
}

func (s *Server) saveToDailyNote(ctx context.Context, _ *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, map[string]any, error) {
	spaceID, err := s.resolveSpaceID(getString(input, "spaceId"))
	if err != nil {
		return nil, nil, err
	}

	text := getString(input, "text")
	if text == "" {
		return nil, nil, fmt.Errorf("text is required")
	}

	noTimeStamp := getBool(input, "noTimeStamp")

	if err := s.client.SaveToDailyNote(ctx, spaceID, text, api.SaveOptions{NoTimeStamp: noTimeStamp}); err != nil {
		return nil, nil, fmt.Errorf("failed to save to daily note: %w", err)
	}

	return nil, map[string]any{
		"message":     "Successfully saved to daily note",
		"spaceId":     spaceID,
		"text":        text,
		"noTimeStamp": noTimeStamp,
	}, nil
}

func (s *Server) search(ctx context.Context, _ *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, map[string]any, error) {
	spaceID, err := s.resolveSpaceID(getString(input, "spaceId"))
	if err != nil {
		return nil, nil, err
	}

	searchTerm := getString(input, "searchTerm")
	if searchTerm == "" {
		return nil, nil, fmt.Errorf("search term is required")
	}

	results, err := s.client.Lookup(ctx, api.LookupRequest{
		SearchTerm: searchTerm,
		SpaceID:    spaceID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("lookup failed: %w", err)
	}

	return nil, map[string]any{
		"spaceId":    spaceID,
		"searchTerm": searchTerm,
		"results":    results,
	}, nil
}

func (s *Server) listSpaces(ctx context.Context, _ *mcp.CallToolRequest, _ map[string]any) (*mcp.CallToolResult, map[string]any, error) {
	spaces, err := s.client.GetSpaces(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get spaces: %w", err)
	}

	return nil, map[string]any{"spaces": spaces}, nil
}

func (s *Server) spaceInfo(ctx context.Context, _ *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, map[string]any, error) {
	spaceID, err := s.resolveSpaceID(getString(input, "spaceId"))
	if err != nil {
		return nil, nil, err
	}

	structures, err := s.client.GetSpaceInfo(ctx, spaceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get space info: %w", err)
	}

	return nil, map[string]any{
		"spaceId":    spaceID,
		"structures": structures,
	}, nil
}

func (s *Server) saveWebLink(ctx context.Context, _ *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, map[string]any, error) {
	spaceID, err := s.resolveSpaceID(getString(input, "spaceId"))
	if err != nil {
		return nil, nil, err
	}

	urlStr := getString(input, "url")
	if urlStr == "" {
		return nil, nil, fmt.Errorf("url is required")
	}

	resp, err := s.client.SaveWebLink(ctx, api.SaveWebLinkRequest{
		SpaceID:              spaceID,
		URL:                  urlStr,
		TitleOverwrite:       getString(input, "title"),
		DescriptionOverwrite: getString(input, "description"),
		Tags:                 getStringSlice(input, "tags"),
		MDText:               getString(input, "mdText"),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save weblink: %w", err)
	}

	return nil, map[string]any{"response": resp}, nil
}

func (s *Server) resolveSpaceID(spaceID string) (string, error) {
	if spaceID != "" {
		return spaceID, nil
	}
	if s.cfg != nil && s.cfg.DefaultSpaceID != "" {
		return s.cfg.DefaultSpaceID, nil
	}
	return "", fmt.Errorf("space ID is required (either via input or config)")
}

// Helper functions for map access
func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok && v != nil {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func getStringSlice(m map[string]any, key string) []string {
	if v, ok := m[key]; ok && v != nil {
		if arr, ok := v.([]any); ok {
			result := make([]string, len(arr))
			for i, item := range arr {
				if s, ok := item.(string); ok {
					result[i] = s
				}
			}
			return result
		}
		if arr, ok := v.([]string); ok {
			return arr
		}
	}
	return nil
}
