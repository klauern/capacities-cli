package mcpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/klauern/capacities-cli/internal/api"
	"github.com/klauern/capacities-cli/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestServerTools(t *testing.T) {
	t.Parallel()

	var (
		mu           sync.Mutex
		dailyRequest map[string]any
		searchQuery  url.Values
		spaceInfoQ   url.Values
		webLinkBody  map[string]any
	)

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Bearer token-123"; got != want {
			t.Errorf("request %s: auth header = %q, want %q", r.URL.Path, got, want)
		}

		switch r.URL.Path {
		case "/save-to-daily-note":
			defer func() { _ = r.Body.Close() }()
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode daily request: %v", err)
			}
			mu.Lock()
			dailyRequest = body
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
		case "/lookup":
			defer func() { _ = r.Body.Close() }()
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode lookup request: %v", err)
			}
			mu.Lock()
			searchQuery = r.URL.Query()
			mu.Unlock()
			if got, want := body["searchTerm"], "alpha"; got != want {
				t.Errorf("lookup searchTerm = %v, want %v", got, want)
			}
			if got, want := body["spaceId"], "space-2"; got != want {
				t.Errorf("lookup spaceId = %v, want %v", got, want)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{
					{"id": "result-1", "structureId": "structure-1", "title": "Alpha"},
				},
			})
		case "/spaces":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"spaces": []map[string]any{
					{"id": "space-1", "title": "Personal", "icon": map[string]any{"val": "📁"}},
				},
			})
		case "/space-info":
			mu.Lock()
			spaceInfoQ = r.URL.Query()
			mu.Unlock()
			if got, want := r.URL.Query().Get("spaceId"), "space-1"; got != want {
				t.Errorf("space-info spaceId = %q, want %q", got, want)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"structures": []map[string]any{
					{"id": "structure-1", "title": "Daily Notes", "collections": []map[string]any{{"title": "Default"}}},
				},
			})
		case "/save-weblink":
			defer func() { _ = r.Body.Close() }()
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode weblink request: %v", err)
			}
			mu.Lock()
			webLinkBody = body
			mu.Unlock()
			if got, want := body["spaceId"], "space-3"; got != want {
				t.Errorf("weblink spaceId = %v, want %v", got, want)
			}
			if got, want := body["url"], "https://example.com"; got != want {
				t.Errorf("weblink url = %v, want %v", got, want)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"spaceId":     "space-3",
				"id":          "link-1",
				"structureId": "weblink",
				"title":       "Example",
				"description": "Example desc",
				"tags":        []string{"alpha", "beta"},
			})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer apiServer.Close()

	cfg := &config.Config{
		Token:          "token-123",
		DefaultSpaceID: "space-1",
	}
	srv, err := NewServer(cfg, api.NewClientWithBaseURL(cfg.Token, apiServer.URL))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("connect server: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "dev"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	defer clientSession.Close()

	callTool := func(name string, args any) *mcp.CallToolResult {
		t.Helper()
		result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil {
			t.Fatalf("call %s: %v", name, err)
		}
		if result.IsError {
			t.Fatalf("call %s returned error: %s", name, result.Content[0].(*mcp.TextContent).Text)
		}
		return result
	}

	daily := callTool("save_to_daily_note", map[string]any{
		"text":        "hello daily",
		"noTimeStamp": true,
	})
	dailyContent := mustStructuredMap(t, daily.StructuredContent)
	if got, want := dailyContent["spaceId"], "space-1"; got != want {
		t.Fatalf("daily spaceId = %v, want %v", got, want)
	}
	if got, want := dailyContent["text"], "hello daily"; got != want {
		t.Fatalf("daily text = %v, want %v", got, want)
	}
	if got, want := dailyContent["noTimeStamp"], true; got != want {
		t.Fatalf("daily noTimeStamp = %v, want %v", got, want)
	}

	search := callTool("search", map[string]any{
		"searchTerm": "alpha",
		"spaceId":    "space-2",
	})
	searchContent := mustStructuredMap(t, search.StructuredContent)
	if got, want := searchContent["spaceId"], "space-2"; got != want {
		t.Fatalf("search spaceId = %v, want %v", got, want)
	}
	results, ok := searchContent["results"].([]any)
	if !ok || len(results) != 1 {
		t.Fatalf("search results = %#v, want 1 item", searchContent["results"])
	}

	spaces := callTool("list_spaces", map[string]any{})
	spacesContent := mustStructuredMap(t, spaces.StructuredContent)
	spaceList, ok := spacesContent["spaces"].([]any)
	if !ok || len(spaceList) != 1 {
		t.Fatalf("spaces = %#v, want 1 item", spacesContent["spaces"])
	}

	info := callTool("space_info", map[string]any{})
	infoContent := mustStructuredMap(t, info.StructuredContent)
	if got, want := infoContent["spaceId"], "space-1"; got != want {
		t.Fatalf("space info spaceId = %v, want %v", got, want)
	}

	link := callTool("save_weblink", map[string]any{
		"url":         "https://example.com",
		"spaceId":     "space-3",
		"title":       "Example",
		"description": "Example desc",
		"tags":        []string{"alpha", "beta"},
		"mdText":      "notes",
	})
	linkContent := mustStructuredMap(t, link.StructuredContent)
	response, ok := linkContent["response"].(map[string]any)
	if !ok {
		t.Fatalf("save weblink response = %#v, want map", linkContent["response"])
	}
	if got, want := response["id"], "link-1"; got != want {
		t.Fatalf("weblink id = %v, want %v", got, want)
	}

	mu.Lock()
	defer mu.Unlock()
	if dailyRequest == nil {
		t.Fatal("expected daily request to be recorded")
	}
	if got, want := dailyRequest["spaceId"], "space-1"; got != want {
		t.Fatalf("daily request spaceId = %v, want %v", got, want)
	}
	if got, want := dailyRequest["mdText"], "hello daily"; got != want {
		t.Fatalf("daily request mdText = %v, want %v", got, want)
	}
	if got, want := dailyRequest["noTimeStamp"], true; got != want {
		t.Fatalf("daily request noTimeStamp = %v, want %v", got, want)
	}
	if searchQuery == nil || searchQuery.Get("spaceId") != "" {
		t.Fatalf("expected lookup query to be empty, got %v", searchQuery)
	}
	if spaceInfoQ == nil || spaceInfoQ.Get("spaceId") != "space-1" {
		t.Fatalf("expected space-info query to use default space, got %v", spaceInfoQ)
	}
	if webLinkBody == nil {
		t.Fatal("expected weblink request to be recorded")
	}
	if got, want := webLinkBody["mdText"], "notes"; got != want {
		t.Fatalf("weblink mdText = %v, want %v", got, want)
	}
}

func TestServerRejectsMissingSpaceID(t *testing.T) {
	t.Parallel()

	srv, err := NewServer(&config.Config{Token: "token-123"}, api.NewClientWithBaseURL("token-123", "http://example.invalid"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("connect server: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "dev"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	defer clientSession.Close()

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "search",
		Arguments: map[string]any{"searchTerm": "alpha"},
	})
	if err != nil {
		t.Fatalf("call search: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tool error when no space id is configured")
	}
	if got := result.Content[0].(*mcp.TextContent).Text; !strings.Contains(got, "space ID is required") {
		t.Fatalf("unexpected error content %q", got)
	}
}

func mustStructuredMap(t *testing.T, v any) map[string]any {
	t.Helper()

	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
	return out
}
