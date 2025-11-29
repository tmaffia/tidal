package tidal

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(WithBaseURL(server.URL))

	return func() {
		server.Close()
	}
}

func TestGetArtist(t *testing.T) {
	teardown := setup()
	defer teardown()

	// Mock JSON response from Tidal API v2
	mockResponse := `{
		"data": {
			"id": "1566",
			"type": "artists",
			"attributes": {
				"name": "Beyoncé",
				"popularity": 0.9582937827192956,
				"externalLinks": [
					{
						"href": "https://tidal.com/browse/artist/1566",
						"meta": {
							"type": "TIDAL_SHARING"
						}
					}
				]
			},
			"relationships": {
				"albums": {
					"links": {
						"self": "/artists/1566/relationships/albums?countryCode=US"
					}
				}
			},
			"links": {
				"self": "/artists/1566?countryCode=US"
			}
		}
	}`

	mux.HandleFunc("/artists/1566", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.URL.Query().Get("countryCode") != "US" {
			t.Errorf("Expected countryCode=US, got %s", r.URL.Query().Get("countryCode"))
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockResponse)
	})

	// Call the GetArtist method
	artist, err := client.GetArtist(context.Background(), "1566")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the response
	if artist.Data.ID != "1566" {
		t.Errorf("Expected ID 1566, got %s", artist.Data.ID)
	}
	if artist.Data.Attributes.Name != "Beyoncé" {
		t.Errorf("Expected Name Beyoncé, got %s", artist.Data.Attributes.Name)
	}
}
