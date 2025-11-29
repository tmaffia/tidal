package tidal

import (
	"context"
	"fmt"
)

// ArtistService handles communication with the artist related methods of the Tidal API.
type ArtistService struct {
	client *Client
}

// Artist represents a Tidal artist.
type Artist struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Picture     string   `json:"picture"`
	Popularity  float64  `json:"popularity"`
	ArtistTypes []string `json:"artistTypes"`
}

type artistResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes Artist `json:"attributes"`
	} `json:"data"`
}

// Get returns an artist by ID.
func (s *ArtistService) Get(ctx context.Context, id string, options ...Option) (*Artist, error) {
	path := fmt.Sprintf("artists/%s", id)
	req, err := s.client.NewRequest("GET", path, nil, options...)
	if err != nil {
		return nil, err
	}

	var resp artistResponse
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	artist := resp.Data.Attributes
	artist.ID = resp.Data.ID // ID is often at the top level in JSON:API
	return &artist, nil
}

// GetMultiple returns a list of artists by their IDs.
func (s *ArtistService) GetMultiple(ctx context.Context, ids []string, options ...Option) ([]*Artist, error) {
	// Placeholder
	return nil, fmt.Errorf("GetMultiple not yet implemented")
}
