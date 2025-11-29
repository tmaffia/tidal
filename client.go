package tidal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultBaseURL = "https://openapi.tidal.com/v2"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

func (c *Client) GetArtist(ctx context.Context, id string) (*ArtistResponse, error) {
	url := fmt.Sprintf("%s/artists/%s?countryCode=US", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var artistResponse ArtistResponse
	if err := json.NewDecoder(resp.Body).Decode(&artistResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &artistResponse, nil
}
