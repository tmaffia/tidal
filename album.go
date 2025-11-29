package tidal

import (
	"context"
	"fmt"
)

// AlbumService handles communication with the album related methods of the Tidal API.
type AlbumService struct {
	client *Client
}

// Album represents a Tidal album.
type Album struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Duration             string `json:"duration"`
	StreamReady          bool   `json:"streamReady"`
	StreamStartDate      string `json:"streamStartDate"`
	AllowStreaming       bool   `json:"allowStreaming"`
	PremiumStreamingOnly bool   `json:"premiumStreamingOnly"`
	NumberOfTracks       int    `json:"numberOfTracks"`
	NumberOfItems        int    `json:"numberOfItems"`
	NumberOfVideos       int    `json:"numberOfVideos"`
	NumberOfVolumes      int    `json:"numberOfVolumes"`
	ReleaseDate          string `json:"releaseDate"`
	Copyright            struct {
		Text string `json:"text"`
	} `json:"copyright"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	Url          string    `json:"url"`
	Cover        string    `json:"cover"`
	VideoCover   string    `json:"videoCover"`
	Explicit     bool      `json:"explicit"`
	Upc          string    `json:"upc"`
	Popularity   float64   `json:"popularity"`
	AudioQuality string    `json:"audioQuality"`
	AudioModes   []string  `json:"audioModes"`
	Artist       *Artist   `json:"artist"`
	Artists      []*Artist `json:"artists"`
}

type albumResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes Album  `json:"attributes"`
	} `json:"data"`
}

// Get returns an album by ID.
func (s *AlbumService) Get(ctx context.Context, id string, options ...Option) (*Album, error) {
	path := fmt.Sprintf("albums/%s", id)
	req, err := s.client.NewRequest("GET", path, nil, options...)
	if err != nil {
		return nil, err
	}

	var resp albumResponse
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	album := resp.Data.Attributes
	album.ID = resp.Data.ID
	return &album, nil
}
