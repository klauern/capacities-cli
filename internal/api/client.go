// Package api provides the HTTP client for interacting with the Capacities.io API.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.capacities.io"

// Error represents a non-200 response from the Capacities API.
type Error struct {
	StatusCode int
	Method     string
	URL        string
	Body       []byte
}

// Error returns a string describing the API error, including the HTTP status code and response body.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	body := string(e.Body)
	if body == "" {
		return fmt.Sprintf("API request failed with status %d", e.StatusCode)
	}
	return fmt.Sprintf("API request failed with status %d: %s", e.StatusCode, body)
}

// Client is an HTTP client for interacting with the Capacities.io API.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Capacities API client with the given API token.
func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doJSON performs an authenticated JSON HTTP request, marshaling the request body and
// unmarshaling the response into responseBody. It returns an *Error for non-2xx responses.
func (c *Client) doJSON(ctx context.Context, method string, path string, query url.Values, requestBody any, responseBody any) error {
	var bodyReader io.Reader
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &Error{
			StatusCode: resp.StatusCode,
			Method:     method,
			URL:        req.URL.String(),
			Body:       body,
		}
	}

	if responseBody == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// SaveOptions contains optional settings for saving to daily note.
type SaveOptions struct {
	NoTimeStamp bool
}

type saveToDailyNoteRequest struct {
	SpaceID     string `json:"spaceId"`
	MDText      string `json:"mdText"`
	Origin      string `json:"origin,omitempty"`
	NoTimeStamp bool   `json:"noTimeStamp,omitempty"`
}

// SaveToDailyNote saves text content to the daily note in the specified space.
func (c *Client) SaveToDailyNote(ctx context.Context, spaceID string, text string, opts SaveOptions) error {
	reqBody := saveToDailyNoteRequest{
		SpaceID:     spaceID,
		MDText:      text,
		Origin:      "commandPalette", // Using commandPalette as origin seems appropriate for a CLI
		NoTimeStamp: opts.NoTimeStamp,
	}

	return c.doJSON(ctx, http.MethodPost, "/save-to-daily-note", nil, reqBody, nil)
}

// Icon represents an icon configuration in Capacities.
type Icon struct {
	Type     string `json:"type"`
	Val      string `json:"val"`
	Color    string `json:"color"`
	ColorHex string `json:"colorHex"`
}

// Space represents a workspace in Capacities.
type Space struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  Icon   `json:"icon"`
}

type spacesResponse struct {
	Spaces []Space `json:"spaces"`
}

// GetSpaces retrieves all spaces accessible with the current API token.
func (c *Client) GetSpaces(ctx context.Context) ([]Space, error) {
	var result spacesResponse
	if err := c.doJSON(ctx, http.MethodGet, "/spaces", nil, nil, &result); err != nil {
		return nil, err
	}

	return result.Spaces, nil
}

// PropertyDefinition describes a property in a Capacities structure.
type PropertyDefinition struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	DataType string `json:"dataType"`
	Name     string `json:"name"`
}

// Collection represents a collection within a structure.
type Collection struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Structure represents a content structure (type) in Capacities.
type Structure struct {
	ID                  string               `json:"id"`
	Title               string               `json:"title"`
	PluralName          string               `json:"pluralName"`
	PropertyDefinitions []PropertyDefinition `json:"propertyDefinitions"`
	LabelColor          string               `json:"labelColor"`
	Collections         []Collection         `json:"collections"`
}

type spaceInfoResponse struct {
	Structures []Structure `json:"structures"`
}

// GetSpaceInfo retrieves detailed information about structures and collections in a space.
func (c *Client) GetSpaceInfo(ctx context.Context, spaceID string) ([]Structure, error) {
	//nolint:revive // API expects "spaceId" (camelCase)
	q := url.Values{}
	q.Add("spaceId", spaceID)

	var result spaceInfoResponse
	if err := c.doJSON(ctx, http.MethodGet, "/space-info", q, nil, &result); err != nil {
		return nil, err
	}

	return result.Structures, nil
}

// LookupRequest contains parameters for looking up content by title in Capacities.
type LookupRequest struct {
	SearchTerm string `json:"searchTerm"`
	SpaceID    string `json:"spaceId"`
}

// LookupResult represents a single result from a lookup operation.
type LookupResult struct {
	ID          string `json:"id"`
	StructureID string `json:"structureId"`
	Title       string `json:"title"`
}

type lookupResponse struct {
	Results []LookupResult `json:"results"`
}

// Lookup searches for content by title in the specified space.
func (c *Client) Lookup(ctx context.Context, req LookupRequest) ([]LookupResult, error) {
	var result lookupResponse
	if err := c.doJSON(ctx, http.MethodPost, "/lookup", nil, req, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// SaveWebLinkRequest contains parameters for saving a web link to Capacities.
type SaveWebLinkRequest struct {
	SpaceID              string   `json:"spaceId"`
	URL                  string   `json:"url"`
	TitleOverwrite       string   `json:"titleOverwrite,omitempty"`
	DescriptionOverwrite string   `json:"descriptionOverwrite,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	MDText               string   `json:"mdText,omitempty"`
}

// SaveWebLinkResponse contains the result of saving a web link.
type SaveWebLinkResponse struct {
	SpaceID     string   `json:"spaceId"`
	ID          string   `json:"id"`
	StructureID string   `json:"structureId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// SaveWebLink saves a web link to the specified space in Capacities.
func (c *Client) SaveWebLink(ctx context.Context, req SaveWebLinkRequest) (*SaveWebLinkResponse, error) {
	var result SaveWebLinkResponse
	if err := c.doJSON(ctx, http.MethodPost, "/save-weblink", nil, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
