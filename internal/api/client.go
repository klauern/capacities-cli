package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.capacities.io"

type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
}

type SaveOptions struct {
	NoTimeStamp bool
}

type saveToDailyNoteRequest struct {
	SpaceID     string `json:"spaceId"`
	MDText      string `json:"mdText"`
	Origin      string `json:"origin,omitempty"`
	NoTimeStamp bool   `json:"noTimeStamp,omitempty"`
}

func (c *Client) SaveToDailyNote(ctx context.Context, spaceID string, text string, opts SaveOptions) error {
	reqBody := saveToDailyNoteRequest{
		SpaceID:     spaceID,
		MDText:      text,
		Origin:      "commandPalette", // Using commandPalette as origin seems appropriate for a CLI
		NoTimeStamp: opts.NoTimeStamp,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/save-to-daily-note", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

type Icon struct {
	Type     string `json:"type"`
	Val      string `json:"val"`
	Color    string `json:"color"`
	ColorHex string `json:"colorHex"`
}

type Space struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  Icon   `json:"icon"`
}

type spacesResponse struct {
	Spaces []Space `json:"spaces"`
}

func (c *Client) GetSpaces(ctx context.Context) ([]Space, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/spaces", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result spacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Spaces, nil
}

type PropertyDefinition struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	DataType string `json:"dataType"`
	Name     string `json:"name"`
}

type Collection struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

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

func (c *Client) GetSpaceInfo(ctx context.Context, spaceID string) ([]Structure, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/space-info", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("spaceid", spaceID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result spaceInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Structures, nil
}

type SearchRequest struct {
	SearchTerm         string   `json:"searchTerm"`
	SpaceIDs           []string `json:"spaceIds"`
	FilterStructureIDs []string `json:"filterStructureIds,omitempty"`
	Mode               string   `json:"mode,omitempty"`
}

type Context struct {
	Field string `json:"field"`
}

type Highlight struct {
	Context  Context  `json:"context"`
	Snippets []string `json:"snippets"`
	Score    float64  `json:"score"`
}

type SearchResult struct {
	ID          string      `json:"id"`
	SpaceID     string      `json:"spaceId"`
	StructureID string      `json:"structureId"`
	Title       string      `json:"title"`
	Highlights  []Highlight `json:"highlights"`
}

type searchResponse struct {
	Results []SearchResult `json:"results"`
}

func (c *Client) Search(ctx context.Context, req SearchRequest) ([]SearchResult, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/search", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Results, nil
}

type SaveWebLinkRequest struct {
	SpaceID              string   `json:"spaceId"`
	URL                  string   `json:"url"`
	TitleOverwrite       string   `json:"titleOverwrite,omitempty"`
	DescriptionOverwrite string   `json:"descriptionOverwrite,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	MDText               string   `json:"mdText,omitempty"`
}

type SaveWebLinkResponse struct {
	SpaceID     string   `json:"spaceId"`
	ID          string   `json:"id"`
	StructureID string   `json:"structureId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (c *Client) SaveWebLink(ctx context.Context, req SaveWebLinkRequest) (*SaveWebLinkResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/save-weblink", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result SaveWebLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
