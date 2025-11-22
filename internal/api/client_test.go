package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SaveToDailyNote(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/save-to-daily-note" {
			t.Errorf("Expected path /save-to-daily-note, got %s", r.URL.Path)
		}

		// Verify headers
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header 'application/json', got %s", r.Header.Get("Content-Type"))
		}

		// Verify body
		var reqBody saveToDailyNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if reqBody.SpaceID != "test-space" {
			t.Errorf("Expected SpaceID 'test-space', got %s", reqBody.SpaceID)
		}
		if reqBody.MDText != "test note" {
			t.Errorf("Expected MDText 'test note', got %s", reqBody.MDText)
		}

		// Send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "note-id"}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	// We need to modify the client to allow overriding the base URL for testing
	// For now, we can temporarily change the constant or make it configurable.
	// A better approach is to add a BaseURL field to the Client struct.

	// Since I cannot change the const in the same package easily without modifying the struct,
	// I will modify the Client struct in client.go to support BaseURL.

	client := NewClient("test-token")
	client.baseURL = server.URL // This requires modifying Client struct to have baseURL field

	if err := client.SaveToDailyNote(context.Background(), "test-space", "test note", SaveOptions{}); err != nil {
		t.Fatalf("SaveToDailyNote failed: %v", err)
	}
}

func TestClient_GetSpaces(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/spaces" {
			t.Errorf("Expected path /spaces, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"spaces": [
				{
					"id": "space-1",
					"title": "My Space",
					"icon": {
						"type": "emoji",
						"val": "🚀",
						"color": "blue",
						"colorHex": "#0000FF"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	spaces, err := client.GetSpaces(context.Background())
	if err != nil {
		t.Fatalf("GetSpaces failed: %v", err)
	}

	if len(spaces) != 1 {
		t.Errorf("Expected 1 space, got %d", len(spaces))
	}
	if spaces[0].ID != "space-1" {
		t.Errorf("Expected space ID 'space-1', got %s", spaces[0].ID)
	}
	if spaces[0].Title != "My Space" {
		t.Errorf("Expected space title 'My Space', got %s", spaces[0].Title)
	}
}

func TestClient_GetSpaceInfo(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/space-info" {
			t.Errorf("Expected path /space-info, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("spaceId") != "space-1" {
			t.Errorf("Expected spaceId query param 'space-1', got %s", r.URL.Query().Get("spaceId"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"structures": [
				{
					"id": "struct-1",
					"title": "Book",
					"pluralName": "Books",
					"propertyDefinitions": [
						{
							"id": "prop-1",
							"type": "text",
							"dataType": "string",
							"name": "Author"
						}
					],
					"labelColor": "blue",
					"collections": [
						{
							"id": "col-1",
							"title": "Fiction"
						}
					]
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	structures, err := client.GetSpaceInfo(context.Background(), "space-1")
	if err != nil {
		t.Fatalf("GetSpaceInfo failed: %v", err)
	}

	if len(structures) != 1 {
		t.Errorf("Expected 1 structure, got %d", len(structures))
	}
	if structures[0].ID != "struct-1" {
		t.Errorf("Expected structure ID 'struct-1', got %s", structures[0].ID)
	}
	if structures[0].Title != "Book" {
		t.Errorf("Expected structure title 'Book', got %s", structures[0].Title)
	}
	if len(structures[0].PropertyDefinitions) != 1 {
		t.Errorf("Expected 1 property definition, got %d", len(structures[0].PropertyDefinitions))
	}
	if structures[0].PropertyDefinitions[0].Name != "Author" {
		t.Errorf("Expected property name 'Author', got %s", structures[0].PropertyDefinitions[0].Name)
	}
	if len(structures[0].Collections) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(structures[0].Collections))
	}
	if structures[0].Collections[0].Title != "Fiction" {
		t.Errorf("Expected collection title 'Fiction', got %s", structures[0].Collections[0].Title)
	}
}

func TestClient_Search(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/search" {
			t.Errorf("Expected path /search, got %s", r.URL.Path)
		}

		var reqBody SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if reqBody.SearchTerm != "AI" {
			t.Errorf("Expected SearchTerm 'AI', got %s", reqBody.SearchTerm)
		}
		if len(reqBody.SpaceIDs) != 1 || reqBody.SpaceIDs[0] != "space-1" {
			t.Errorf("Expected SpaceIDs ['space-1'], got %v", reqBody.SpaceIDs)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [
				{
					"id": "res-1",
					"spaceId": "space-1",
					"structureId": "struct-1",
					"title": "AI Research",
					"highlights": [
						{
							"context": {
								"field": "title"
							},
							"snippets": ["AI"],
							"score": 1.0
						}
					]
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	results, err := client.Search(context.Background(), SearchRequest{
		SearchTerm: "AI",
		SpaceIDs:   []string{"space-1"},
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].ID != "res-1" {
		t.Errorf("Expected result ID 'res-1', got %s", results[0].ID)
	}
	if results[0].Title != "AI Research" {
		t.Errorf("Expected result title 'AI Research', got %s", results[0].Title)
	}
}

func TestClient_SaveWebLink(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/save-weblink" {
			t.Errorf("Expected path /save-weblink, got %s", r.URL.Path)
		}

		var reqBody SaveWebLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if reqBody.URL != "https://example.com" {
			t.Errorf("Expected URL 'https://example.com', got %s", reqBody.URL)
		}
		if reqBody.SpaceID != "space-1" {
			t.Errorf("Expected SpaceID 'space-1', got %s", reqBody.SpaceID)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"spaceId": "space-1",
			"id": "link-1",
			"structureId": "struct-1",
			"title": "Example",
			"description": "An example",
			"tags": ["test"]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	resp, err := client.SaveWebLink(context.Background(), SaveWebLinkRequest{
		SpaceID: "space-1",
		URL:     "https://example.com",
	})
	if err != nil {
		t.Fatalf("SaveWebLink failed: %v", err)
	}

	if resp.ID != "link-1" {
		t.Errorf("Expected ID 'link-1', got %s", resp.ID)
	}
	if resp.Title != "Example" {
		t.Errorf("Expected title 'Example', got %s", resp.Title)
	}
}
