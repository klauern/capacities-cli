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

type saveToDailyNoteInput struct {
	Text        string `json:"text" jsonschema:"text to append to today's daily note"`
	SpaceID     string `json:"spaceId,omitempty" jsonschema:"override the configured default space id"`
	NoTimeStamp bool   `json:"noTimeStamp,omitempty" jsonschema:"do not prefix the note with a timestamp"`
}

type saveToDailyNoteOutput struct {
	Message     string `json:"message"`
	SpaceID     string `json:"spaceId"`
	Text        string `json:"text"`
	NoTimeStamp bool   `json:"noTimeStamp"`
}

func (s *Server) saveToDailyNote(ctx context.Context, _ *mcp.CallToolRequest, input saveToDailyNoteInput) (*mcp.CallToolResult, saveToDailyNoteOutput, error) {
	spaceID, err := s.resolveSpaceID(input.SpaceID)
	if err != nil {
		return nil, saveToDailyNoteOutput{}, err
	}

	if input.Text == "" {
		return nil, saveToDailyNoteOutput{}, fmt.Errorf("text is required")
	}

	if err := s.client.SaveToDailyNote(ctx, spaceID, input.Text, api.SaveOptions{NoTimeStamp: input.NoTimeStamp}); err != nil {
		return nil, saveToDailyNoteOutput{}, fmt.Errorf("failed to save to daily note: %w", err)
	}

	return nil, saveToDailyNoteOutput{
		Message:     "Successfully saved to daily note",
		SpaceID:     spaceID,
		Text:        input.Text,
		NoTimeStamp: input.NoTimeStamp,
	}, nil
}

type searchInput struct {
	SearchTerm string `json:"searchTerm" jsonschema:"search term to look up"`
	SpaceID    string `json:"spaceId,omitempty" jsonschema:"override the configured default space id"`
}

type searchOutput struct {
	SpaceID    string             `json:"spaceId"`
	SearchTerm string             `json:"searchTerm"`
	Results    []api.LookupResult `json:"results"`
}

func (s *Server) search(ctx context.Context, _ *mcp.CallToolRequest, input searchInput) (*mcp.CallToolResult, searchOutput, error) {
	spaceID, err := s.resolveSpaceID(input.SpaceID)
	if err != nil {
		return nil, searchOutput{}, err
	}
	if input.SearchTerm == "" {
		return nil, searchOutput{}, fmt.Errorf("search term is required")
	}

	results, err := s.client.Lookup(ctx, api.LookupRequest{
		SearchTerm: input.SearchTerm,
		SpaceID:    spaceID,
	})
	if err != nil {
		return nil, searchOutput{}, fmt.Errorf("lookup failed: %w", err)
	}

	return nil, searchOutput{
		SpaceID:    spaceID,
		SearchTerm: input.SearchTerm,
		Results:    results,
	}, nil
}

type listSpacesInput struct{}

type listSpacesOutput struct {
	Spaces []api.Space `json:"spaces"`
}

func (s *Server) listSpaces(ctx context.Context, _ *mcp.CallToolRequest, _ listSpacesInput) (*mcp.CallToolResult, listSpacesOutput, error) {
	spaces, err := s.client.GetSpaces(ctx)
	if err != nil {
		return nil, listSpacesOutput{}, fmt.Errorf("failed to get spaces: %w", err)
	}

	return nil, listSpacesOutput{Spaces: spaces}, nil
}

type spaceInfoInput struct {
	SpaceID string `json:"spaceId,omitempty" jsonschema:"override the configured default space id"`
}

type spaceInfoOutput struct {
	SpaceID    string          `json:"spaceId"`
	Structures []api.Structure `json:"structures"`
}

func (s *Server) spaceInfo(ctx context.Context, _ *mcp.CallToolRequest, input spaceInfoInput) (*mcp.CallToolResult, spaceInfoOutput, error) {
	spaceID, err := s.resolveSpaceID(input.SpaceID)
	if err != nil {
		return nil, spaceInfoOutput{}, err
	}

	structures, err := s.client.GetSpaceInfo(ctx, spaceID)
	if err != nil {
		return nil, spaceInfoOutput{}, fmt.Errorf("failed to get space info: %w", err)
	}

	return nil, spaceInfoOutput{
		SpaceID:    spaceID,
		Structures: structures,
	}, nil
}

type saveWebLinkInput struct {
	URL                  string   `json:"url" jsonschema:"the web page url to save"`
	SpaceID              string   `json:"spaceId,omitempty" jsonschema:"override the configured default space id"`
	TitleOverwrite       string   `json:"title,omitempty" jsonschema:"optional title override"`
	DescriptionOverwrite string   `json:"description,omitempty" jsonschema:"optional description override"`
	Tags                 []string `json:"tags,omitempty" jsonschema:"optional tags to attach"`
	MDText               string   `json:"mdText,omitempty" jsonschema:"optional markdown note to attach"`
}

type saveWebLinkOutput struct {
	Response api.SaveWebLinkResponse `json:"response"`
}

func (s *Server) saveWebLink(ctx context.Context, _ *mcp.CallToolRequest, input saveWebLinkInput) (*mcp.CallToolResult, saveWebLinkOutput, error) {
	spaceID, err := s.resolveSpaceID(input.SpaceID)
	if err != nil {
		return nil, saveWebLinkOutput{}, err
	}
	if input.URL == "" {
		return nil, saveWebLinkOutput{}, fmt.Errorf("url is required")
	}

	resp, err := s.client.SaveWebLink(ctx, api.SaveWebLinkRequest{
		SpaceID:              spaceID,
		URL:                  input.URL,
		TitleOverwrite:       input.TitleOverwrite,
		DescriptionOverwrite: input.DescriptionOverwrite,
		Tags:                 input.Tags,
		MDText:               input.MDText,
	})
	if err != nil {
		return nil, saveWebLinkOutput{}, fmt.Errorf("failed to save weblink: %w", err)
	}

	return nil, saveWebLinkOutput{Response: *resp}, nil
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
