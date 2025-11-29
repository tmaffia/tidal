package tidal

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// Authenticator handles the OAuth2 Authorization Code Flow with PKCE.
type Authenticator struct {
	config *oauth2.Config
}

// AuthenticatorOption is a function that configures the Authenticator.
type AuthenticatorOption func(*authenticatorConfig)

type authenticatorConfig struct {
	tokenURL  string
	loginURL  string
	authStyle oauth2.AuthStyle
	scopes    []string
}

// NewAuthenticator creates a new Authenticator with the provided credentials.
func NewAuthenticator(clientID, redirectURL string, opts ...AuthenticatorOption) *Authenticator {
	cfg := &authenticatorConfig{
		tokenURL:  defaultAuthURL,
		loginURL:  "https://login.tidal.com",
		authStyle: oauth2.AuthStyleInParams,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &Authenticator{
		config: &oauth2.Config{
			ClientID:    clientID,
			RedirectURL: redirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:   fmt.Sprintf("%s/authorize", cfg.loginURL),
				TokenURL:  fmt.Sprintf("%s/oauth2/token", cfg.tokenURL),
				AuthStyle: cfg.authStyle,
			},
			Scopes: cfg.scopes,
		},
	}
}

func WithAuthenticatorBaseURL(url string) AuthenticatorOption {
	return func(c *authenticatorConfig) {
		c.tokenURL = url
	}
}

// WithScopes configures the scopes to be requested during authorization.
// Common scopes include:
// - "user.read": Read access to user's account
// - "user.write": Write access to user's account
// - "playlists.read": Read access to playlists
// - "playlists.write": Write access to playlists
//
// For a complete list of scopes required for specific endpoints, please refer to the
// Tidal API Reference documentation: https://developer.tidal.com/documentation/api-sdk/api-sdk-authorization
func WithScopes(scopes ...string) AuthenticatorOption {
	return func(c *authenticatorConfig) {
		c.scopes = scopes
	}
}

// AuthCodeURL returns the URL to redirect the user to for authentication,
// and the code verifier that must be used in the subsequent Exchange call.
// The state parameter is used to prevent CSRF attacks.
func (a *Authenticator) AuthCodeURL(state string) (string, string) {
	verifier := oauth2.GenerateVerifier()
	url := a.config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	return url, verifier
}

// Exchange converts an authorization code into a token.
// The verifier must be the same one returned by AuthCodeURL.
func (a *Authenticator) Exchange(ctx context.Context, code, verifier string) (*oauth2.Token, error) {
	return a.config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
}

// TokenSource returns a TokenSource that will return the provided token
// and automatically refresh it as necessary.
func (a *Authenticator) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return a.config.TokenSource(ctx, token)
}
