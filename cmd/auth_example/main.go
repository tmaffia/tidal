package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/tmaffia/tidal"
)

// stateStore stores the state and verifier for the auth flow.
// In a real application, this should be a more robust storage (e.g., session, database).
var (
	stateStore = make(map[string]string)
	storeMu    sync.RWMutex
)

func main() {
	// Load .env file
	env, err := loadEnv(".env")
	if err != nil {
		// Fallback to reading from root if running from cmd/auth_example
		env, err = loadEnv("../../.env")
		if err != nil {
			log.Printf("Warning: could not load .env file: %v", err)
		}
	}

	clientID := getEnv(env, "TIDAL_CLIENT_ID")
	redirectURL := "http://localhost:8080/callback"

	if clientID == "" {
		log.Fatal("TIDAL_CLIENT_ID must be set")
	}

	// Client Secret is optional for PKCE
	authenticator := tidal.NewAuthenticator(clientID, redirectURL, tidal.WithScopes("user.read"))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		state := generateRandomString(32)
		authURL, verifier := authenticator.AuthCodeURL(state)

		storeMu.Lock()
		stateStore[state] = verifier
		storeMu.Unlock()

		fmt.Printf("Redirecting to: %s\n", authURL)
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		fmt.Printf("Callback received. State: %s, Code: %s\n", state, code)

		if state == "" || code == "" {
			http.Error(w, "Missing state or code", http.StatusBadRequest)
			return
		}

		storeMu.RLock()
		verifier, ok := stateStore[state]
		storeMu.RUnlock()

		if !ok {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Clean up state
		storeMu.Lock()
		delete(stateStore, state)
		storeMu.Unlock()

		token, err := authenticator.Exchange(r.Context(), code, verifier)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to exchange token: %v", err)
			fmt.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		fmt.Printf("Token received: Type=%s, AccessToken (len)=%d, RefreshToken (len)=%d, Expiry=%s\n",
			token.TokenType, len(token.AccessToken), len(token.RefreshToken), token.Expiry)
		if len(token.AccessToken) > 0 {
			fmt.Printf("AccessToken prefix: %s...\n", token.AccessToken[:5])
		}

		// Create a client with the token source
		ts := authenticator.TokenSource(r.Context(), token)
		client := tidal.NewClient(tidal.WithTokenSource(ts))

		// Verify with a sample API call (e.g., GetArtist - Beyoncé)
		artist, err := client.GetArtist(r.Context(), "1566")
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get artist: %v", err)
			fmt.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Authentication successful!",
			"token":   token,
			"artist":  artist,
		}); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	fmt.Println("Starting server on :8080")
	fmt.Printf("Redirect URI: %s\n", redirectURL)
	fmt.Println("Go to http://localhost:8080/login to start the auth flow")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func loadEnv(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env, scanner.Err()
}

func getEnv(env map[string]string, key string) string {
	if val, ok := env[key]; ok {
		return val
	}
	return os.Getenv(key)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
