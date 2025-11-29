package tidal

type ArtistResponse struct {
	Data ArtistData `json:"data"`
}

type ArtistData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes ArtistAttributes `json:"attributes"`
}

type ArtistAttributes struct {
	Name          string         `json:"name"`
	Popularity    float64        `json:"popularity"`
	ExternalLinks []ExternalLink `json:"externalLinks"`
}

type ExternalLink struct {
	Href string `json:"href"`
}
