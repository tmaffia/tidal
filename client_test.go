package tidal

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestGetArtist(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/artists/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method 'GET', got %s", r.Method)
		}
		fmt.Fprint(w, `{"data": {"id": "123", "type": "artists", "attributes": {"name": "Test Artist"}}}`)
	})

	artist, err := client.Artists.Get(context.Background(), "123")
	if err != nil {
		t.Errorf("Artists.Get returned error: %v", err)
	}

	want := &Artist{ID: "123", Name: "Test Artist"}
	if !reflect.DeepEqual(artist, want) {
		t.Errorf("Artists.Get returned %+v, want %+v", artist, want)
	}
}

func TestGetAlbum(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/albums/456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method 'GET', got %s", r.Method)
		}
		fmt.Fprint(w, `{"data": {"id": "456", "type": "albums", "attributes": {"title": "Test Album", "duration": "PT3M"}}}`)
	})

	album, err := client.Albums.Get(context.Background(), "456")
	if err != nil {
		t.Errorf("Albums.Get returned error: %v", err)
	}

	want := &Album{ID: "456", Title: "Test Album", Duration: "PT3M"}
	if !reflect.DeepEqual(album, want) {
		t.Errorf("Albums.Get returned %+v, want %+v", album, want)
	}
}

func TestGetSong(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/tracks/789", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method 'GET', got %s", r.Method)
		}
		fmt.Fprint(w, `{"data": {"id": "789", "type": "tracks", "attributes": {"title": "Test Song", "duration": "PT3M"}}}`)
	})

	song, err := client.Songs.Get(context.Background(), "789")
	if err != nil {
		t.Errorf("Songs.Get returned error: %v", err)
	}

	want := &Song{ID: "789", Title: "Test Song", Duration: "PT3M"}
	if !reflect.DeepEqual(song, want) {
		t.Errorf("Songs.Get returned %+v, want %+v", song, want)
	}
}

func TestGetUserFollowedArtists(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/user123/favorites/artists", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method 'GET', got %s", r.Method)
		}
		fmt.Fprint(w, `{"data": [{"id": "123", "type": "artists", "attributes": {"name": "Test Artist"}}]}`)
	})

	artists, err := client.Users.GetFollowedArtists(context.Background(), "user123")
	if err != nil {
		t.Errorf("Users.GetFollowedArtists returned error: %v", err)
	}

	want := []*Artist{{ID: "123", Name: "Test Artist"}}
	if !reflect.DeepEqual(artists, want) {
		t.Errorf("Users.GetFollowedArtists returned %+v, want %+v", artists, want)
	}
}
