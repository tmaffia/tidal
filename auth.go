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
	baseURL string
}

// NewAuthenticator creates a new Authenticator with the provided credentials.
func NewAuthenticator(clientID, clientSecret, redirectURL string, opts ...AuthenticatorOption) *Authenticator {
	cfg := &authenticatorConfig{
		baseURL: defaultAuthURL,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &Authenticator{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:   fmt.Sprintf("%s/oauth2/authorize", cfg.baseURL),
				TokenURL:  fmt.Sprintf("%s/oauth2/token", cfg.baseURL),
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}

// WithAuthenticatorBaseURL sets the base URL for the authentication endpoints.
func WithAuthenticatorBaseURL(url string) AuthenticatorOption {
	return func(c *authenticatorConfig) {
		c.baseURL = url
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
