package tidal

import (
	"context"
	"fmt"
)

// SongService handles communication with the song/track related methods of the Tidal API.
type SongService struct {
	client *Client
}

// Song represents a Tidal song/track.
type Song struct {
	ID                   string  `json:"id"`
	Title                string  `json:"title"`
	Duration             string  `json:"duration"`
	ReplayGain           float64 `json:"replayGain"`
	Peak                 float64 `json:"peak"`
	AllowStreaming       bool    `json:"allowStreaming"`
	StreamReady          bool    `json:"streamReady"`
	StreamStartDate      string  `json:"streamStartDate"`
	PremiumStreamingOnly bool    `json:"premiumStreamingOnly"`
	TrackNumber          int     `json:"trackNumber"`
	VolumeNumber         int     `json:"volumeNumber"`
	Version              string  `json:"version"`
	Popularity           float64 `json:"popularity"`
	Copyright            struct {
		Text string `json:"text"`
	} `json:"copyright"`
	URL          string    `json:"url"`
	ISRC         string    `json:"isrc"`
	Editable     bool      `json:"editable"`
	Explicit     bool      `json:"explicit"`
	AudioQuality string    `json:"audioQuality"`
	AudioModes   []string  `json:"audioModes"`
	Artist       *Artist   `json:"artist"`
	Artists      []*Artist `json:"artists"`
	Album        *Album    `json:"album"`
}

type songResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes Song   `json:"attributes"`
	} `json:"data"`
}

// Get returns a song by ID.
func (s *SongService) Get(ctx context.Context, id string, options ...Option) (*Song, error) {
	path := fmt.Sprintf("tracks/%s", id)
	req, err := s.client.NewRequest("GET", path, nil, options...)
	if err != nil {
		return nil, err
	}

	var resp songResponse
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	song := resp.Data.Attributes
	song.ID = resp.Data.ID
	return &song, nil
}
