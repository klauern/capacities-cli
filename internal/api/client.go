package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.capacities.io"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/save-to-daily-note", bytes.NewBuffer(data))
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
