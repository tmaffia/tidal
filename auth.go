package tidal

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	TokenURL = "https://auth.tidal.com/v1/oauth2/token"
)

// Config holds the configuration for the Tidal API.
type Config struct {
	ClientID     string
	ClientSecret string
	TokenURL     string // Optional: Defaults to Tidal production URL
}

// NewClientCredentialsClient returns an http.Client authenticated via the Client Credentials flow.
// This is useful for accessing public data (Artists, Albums, etc.).
func NewClientCredentialsClient(ctx context.Context, config Config) *http.Client {
	tokenURL := config.TokenURL
	if tokenURL == "" {
		tokenURL = TokenURL
	}
	conf := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     tokenURL,
	}
	return conf.Client(ctx)
}

// NewTokenSourceClient returns an http.Client authenticated using the provided TokenSource.
// This is useful if you already have a token (e.g., for a specific user).
func NewTokenSourceClient(ctx context.Context, ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(ctx, ts)
}
