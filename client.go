package tidal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	defaultBaseURL = "https://openapi.tidal.com/v2"
	defaultAuthURL = "https://auth.tidal.com/v1"
)

type Client struct {
	baseURL      string
	authURL      string
	httpClient   *http.Client
	clientID     string
	clientSecret string
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		authURL:    defaultAuthURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.clientID != "" && c.clientSecret != "" {
		config := &clientcredentials.Config{
			ClientID:     c.clientID,
			ClientSecret: c.clientSecret,
			TokenURL:     fmt.Sprintf("%s/oauth2/token", c.authURL),
			AuthStyle:    oauth2.AuthStyleInParams,
		}
		c.httpClient = config.Client(context.Background())
	}

	return c
}

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithAuthURL(url string) ClientOption {
	return func(c *Client) {
		c.authURL = url
	}
}

func WithClientCredentials(clientID, clientSecret string) ClientOption {
	return func(c *Client) {
		c.clientID = clientID
		c.clientSecret = clientSecret
	}
}

func (c *Client) GetArtist(ctx context.Context, id string) (*ArtistResponse, error) {
	reqURL, err := url.Parse(fmt.Sprintf("%s/artists/%s", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	q := url.Values{}
	q.Set("countryCode", "US")
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
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
