package tidal

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientCredentialsClient(t *testing.T) {
	// Mock the Tidal Token Endpoint AND API Endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			if r.Method != "POST" {
				t.Errorf("Expected method 'POST' for token, got %s", r.Method)
			}
			err := r.ParseForm()
			if err != nil {
				t.Errorf("ParseForm failed: %v", err)
			}
			if r.Form.Get("grant_type") != "client_credentials" {
				t.Errorf("Expected grant_type 'client_credentials', got %s", r.Form.Get("grant_type"))
			}
			// x/oauth2 uses Basic Auth by default for client_credentials.
			clientID := r.Form.Get("client_id")
			clientSecret := r.Form.Get("client_secret")

			if clientID == "" || clientSecret == "" {
				user, pass, ok := r.BasicAuth()
				if !ok {
					t.Error("Expected Basic Auth or form params for client_id/secret")
				}
				clientID = user
				clientSecret = pass
			}

			if clientID != "test-client-id" {
				t.Errorf("Expected client_id 'test-client-id', got %s", clientID)
			}
			if clientSecret != "test-client-secret" {
				t.Errorf("Expected client_secret 'test-client-secret', got %s", clientSecret)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"access_token": "mock-access-token", "token_type": "Bearer", "expires_in": 3600}`)
			return
		}

		if r.URL.Path == "/some-api-endpoint" {
			// Verify the Authorization header
			auth := r.Header.Get("Authorization")
			if auth != "Bearer mock-access-token" {
				t.Errorf("Expected Authorization 'Bearer mock-access-token', got %s", auth)
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		t.Errorf("Unexpected request to %s", r.URL.Path)
	}))
	defer ts.Close()

	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     ts.URL + "/oauth2/token",
	}

	client := NewClientCredentialsClient(context.Background(), config)

	// Verify the client works by making a dummy request (which will trigger the token fetch)
	resp, err := client.Get(ts.URL + "/some-api-endpoint")
	if err != nil {
		t.Fatalf("Client request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
