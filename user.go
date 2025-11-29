package tidal

import (
	"context"
	"fmt"
)

// UserService handles communication with the user related methods of the Tidal API.
type UserService struct {
	client *Client
}

// UserFollowedArtistsResponse represents the response for user's followed artists.
type UserFollowedArtistsResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes Artist `json:"attributes"`
	} `json:"data"`
	Links struct {
		Next string `json:"next"`
		Prev string `json:"prev"`
	} `json:"links"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

// GetFollowedArtists returns the artists followed by the user.
func (s *UserService) GetFollowedArtists(ctx context.Context, userID string, options ...Option) ([]*Artist, error) {
	path := fmt.Sprintf("users/%s/favorites/artists", userID)
	req, err := s.client.NewRequest("GET", path, nil, options...)
	if err != nil {
		return nil, err
	}

	var response UserFollowedArtistsResponse
	_, err = s.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	var artists []*Artist
	for _, item := range response.Data {
		artist := item.Attributes
		artist.ID = item.ID
		artists = append(artists, &artist)
	}

	return artists, nil
}
