package tidal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestClientCredentials(t *testing.T) {
	// Mock response
	mockTokenResponse := `{
		"access_token": "mock_access_token",
		"token_type": "Bearer",
		"expires_in": 3600
	}`

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/oauth2/token" {
			t.Errorf("Expected path /oauth2/token, got %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}
		if r.FormValue("client_id") != "test_client_id" {
			t.Errorf("Expected client_id=test_client_id, got %s", r.FormValue("client_id"))
		}
		if r.FormValue("client_secret") != "test_client_secret" {
			t.Errorf("Expected client_secret=test_client_secret, got %s", r.FormValue("client_secret"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockTokenResponse)
	}))
	defer server.Close()

	// Create client with mock auth URL
	client := NewClient(WithAuthURL(server.URL))

	// Test RequestClientCredentials
	token, err := client.RequestClientCredentials(context.Background(), "test_client_id", "test_client_secret")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token.AccessToken != "mock_access_token" {
		t.Errorf("Expected access token mock_access_token, got %s", token.AccessToken)
	}
	if token.TokenType != "Bearer" {
		t.Errorf("Expected token type Bearer, got %s", token.TokenType)
	}
	if token.ExpiresIn != 3600 {
		t.Errorf("Expected expires in 3600, got %d", token.ExpiresIn)
	}
}

func TestRequestPKCEToken(t *testing.T) {
	// Mock response
	mockTokenResponse := `{
		"access_token": "mock_pkce_access_token",
		"token_type": "Bearer",
		"expires_in": 3600,
		"refresh_token": "mock_refresh_token",
		"scope": "r_usr w_usr"
	}`

	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/oauth2/token" {
			t.Errorf("Expected path /oauth2/token, got %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.FormValue("grant_type") != "authorization_code" {
			t.Errorf("Expected grant_type=authorization_code, got %s", r.FormValue("grant_type"))
		}
		if r.FormValue("code") != "test_auth_code" {
			t.Errorf("Expected code=test_auth_code, got %s", r.FormValue("code"))
		}
		if r.FormValue("client_id") != "test_client_id" {
			t.Errorf("Expected client_id=test_client_id, got %s", r.FormValue("client_id"))
		}
		if r.FormValue("code_verifier") != "test_code_verifier" {
			t.Errorf("Expected code_verifier=test_code_verifier, got %s", r.FormValue("code_verifier"))
		}
		// redirect_uri is optional but good to test if we support it, let's assume we do or just skip for now if not in signature
		// For now, let's keep it simple

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockTokenResponse)
	}))
	defer server.Close()

	// Create client with mock auth URL
	client := NewClient(WithAuthURL(server.URL))

	// Test RequestPKCEToken
	token, err := client.RequestPKCEToken(context.Background(), "test_auth_code", "test_code_verifier", "test_client_id")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token.AccessToken != "mock_pkce_access_token" {
		t.Errorf("Expected access token mock_pkce_access_token, got %s", token.AccessToken)
	}
	if token.RefreshToken != "mock_refresh_token" {
		t.Errorf("Expected refresh token mock_refresh_token, got %s", token.RefreshToken)
	}
	if token.Scope != "r_usr w_usr" {
		t.Errorf("Expected scope r_usr w_usr, got %s", token.Scope)
	}
}

func TestRequestClientCredentials_Failures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("client_id") == "bad_id" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error": "invalid_client", "error_description": "Invalid client ID"}`)
			return
		}
		if r.FormValue("client_secret") == "bad_secret" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error": "invalid_client", "error_description": "Invalid client secret"}`)
			return
		}
	}))
	defer server.Close()

	client := NewClient(WithAuthURL(server.URL))

	// Test Invalid ID
	_, err := client.RequestClientCredentials(context.Background(), "bad_id", "secret")
	if err == nil {
		t.Error("Expected error for invalid client ID, got nil")
	}
	var authErr *AuthError
	if errors.As(err, &authErr) {
		if authErr.ErrorType != "invalid_client" {
			t.Errorf("Expected error type invalid_client, got %s", authErr.ErrorType)
		}
	} else {
		t.Errorf("Expected AuthError, got %T: %v", err, err)
	}

	// Test Invalid Secret
	_, err = client.RequestClientCredentials(context.Background(), "id", "bad_secret")
	if err == nil {
		t.Error("Expected error for invalid client secret, got nil")
	}
	if errors.As(err, &authErr) {
		if authErr.ErrorType != "invalid_client" {
			t.Errorf("Expected error type invalid_client, got %s", authErr.ErrorType)
		}
	} else {
		t.Errorf("Expected AuthError, got %T: %v", err, err)
	}
}

func TestRequestPKCEToken_Failures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("code") == "bad_code" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "invalid_grant", "error_description": "Invalid code"}`)
			return
		}
		if r.FormValue("code_verifier") == "bad_verifier" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "invalid_grant", "error_description": "Invalid code verifier"}`)
			return
		}
	}))
	defer server.Close()

	client := NewClient(WithAuthURL(server.URL))

	// Test Invalid Code
	_, err := client.RequestPKCEToken(context.Background(), "bad_code", "verifier", "id")
	if err == nil {
		t.Error("Expected error for invalid code, got nil")
	}
	var authErr *AuthError
	if errors.As(err, &authErr) {
		if authErr.ErrorType != "invalid_grant" {
			t.Errorf("Expected error type invalid_grant, got %s", authErr.ErrorType)
		}
	} else {
		t.Errorf("Expected AuthError, got %T: %v", err, err)
	}

	// Test Invalid Verifier
	_, err = client.RequestPKCEToken(context.Background(), "code", "bad_verifier", "id")
	if err == nil {
		t.Error("Expected error for invalid verifier, got nil")
	}
	if errors.As(err, &authErr) {
		if authErr.ErrorType != "invalid_grant" {
			t.Errorf("Expected error type invalid_grant, got %s", authErr.ErrorType)
		}
	} else {
		t.Errorf("Expected AuthError, got %T: %v", err, err)
	}
}
