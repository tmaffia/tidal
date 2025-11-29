package tidal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestAuthenticator_AuthCodeURL(t *testing.T) {
	auth := NewAuthenticator("client-id", "http://localhost:8080/callback")
	url, verifier := auth.AuthCodeURL("state-token")

	if verifier == "" {
		t.Error("expected verifier to be non-empty")
	}

	if !strings.Contains(url, "client_id=client-id") {
		t.Error("expected url to contain client_id")
	}
	if !strings.Contains(url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback") {
		t.Error("expected url to contain encoded redirect_uri")
	}
	if !strings.Contains(url, "response_type=code") {
		t.Error("expected url to contain response_type=code")
	}
	if !strings.Contains(url, "code_challenge=") {
		t.Error("expected url to contain code_challenge")
	}
	if !strings.Contains(url, "code_challenge_method=S256") {
		t.Error("expected url to contain code_challenge_method=S256")
	}
}

func TestAuthenticator_Exchange(t *testing.T) {
	// Mock the token endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/oauth2/token" {
			t.Errorf("expected path /oauth2/token, got %s", r.URL.Path)
		}

		err := r.ParseForm()
		if err != nil {
			t.Fatal(err)
		}

		// Ensure NO Basic Auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			t.Errorf("expected no Authorization header, got %s", authHeader)
		}

		// Ensure client_id is in the body
		if r.Form.Get("client_id") != "client-id" {
			t.Errorf("expected client_id=client-id, got %s", r.Form.Get("client_id"))
		}

		if r.Form.Get("grant_type") != "authorization_code" {
			t.Errorf("expected grant_type=authorization_code, got %s", r.Form.Get("grant_type"))
		}
		if r.Form.Get("code") != "auth-code" {
			t.Errorf("expected code=auth-code, got %s", r.Form.Get("code"))
		}
		if r.Form.Get("code_verifier") != "verifier" {
			t.Errorf("expected code_verifier=verifier, got %s", r.Form.Get("code_verifier"))
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{
			"access_token": "mock-access-token",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "mock-refresh-token"
		}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer ts.Close()

	// Override the default auth URL for testing using the new functional option
	auth := NewAuthenticator("client-id", "http://localhost:8080/callback", WithAuthenticatorBaseURL(ts.URL))

	token, err := auth.Exchange(context.Background(), "auth-code", "verifier")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.AccessToken != "mock-access-token" {
		t.Errorf("expected access token mock-access-token, got %s", token.AccessToken)
	}
	if token.RefreshToken != "mock-refresh-token" {
		t.Errorf("expected refresh token mock-refresh-token, got %s", token.RefreshToken)
	}
}

func TestWithTokenSource(t *testing.T) {
	// Create a static token source
	token := &oauth2.Token{
		AccessToken: "static-token",
		TokenType:   "Bearer",
	}
	ts := oauth2.StaticTokenSource(token)

	// client := NewClient(WithTokenSource(ts)) // Removed ineffectual assignment

	// We can't easily inspect the internal http client's transport,
	// but we can verify that making a request uses the token.

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer static-token" {
			t.Errorf("expected Authorization header 'Bearer static-token', got '%s'", authHeader)
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{}`)); err != nil { // Minimal valid JSON
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Temporarily override base URL
	client := NewClient(WithTokenSource(ts), WithBaseURL(server.URL))

	// Make a request (GetArtist calls the base URL)
	_, _ = client.GetArtist(context.Background(), "123")
}
