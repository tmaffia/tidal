package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmaffia/tidal"
)

func main() {
	clientID := os.Getenv("TIDAL_CLIENT_ID")
	clientSecret := os.Getenv("TIDAL_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set TIDAL_CLIENT_ID and TIDAL_CLIENT_SECRET environment variables")
	}

	ctx := context.Background()
	config := tidal.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Create an authenticated HTTP client
	httpClient := tidal.NewClientCredentialsClient(ctx, config)

	// Create the Tidal client
	client := tidal.NewClient(httpClient)

	// Example: Get Artist (Michael Martin Murphey, ID: 1000)
	fmt.Println("Fetching Artist (Michael Martin Murphey)...")
	artist, err := client.Artists.Get(ctx, "1000", tidal.WithCountry("US"))
	if err != nil {
		log.Fatalf("Error fetching artist: %v", err)
	}
	fmt.Printf("Artist: %s (ID: %s)\n", artist.Name, artist.ID)

	// Example: Get Album (ID: 459833113)
	fmt.Println("\nFetching Album...")
	album, err := client.Albums.Get(ctx, "459833113", tidal.WithCountry("US"))
	if err != nil {
		log.Fatalf("Error fetching album: %v", err)
	}
	fmt.Printf("Album: %s (ID: %s, Items: %d)\n", album.Title, album.ID, album.NumberOfItems)
}
