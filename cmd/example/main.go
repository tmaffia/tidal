package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmaffia/tidal"
)

func main() {
	// Load .env file
	env, err := loadEnv(".env")
	if err != nil {
		// Fallback to reading from root if running from cmd/example
		env, err = loadEnv("../../.env")
		if err != nil {
			log.Printf("Warning: could not load .env file: %v", err)
		}
	}

	clientID := getEnv(env, "TIDAL_CLIENT_ID")
	clientSecret := getEnv(env, "TIDAL_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("TIDAL_CLIENT_ID and TIDAL_CLIENT_SECRET must be set")
	}

	client := tidal.NewClient(tidal.WithClientCredentials(clientID, clientSecret))

	requestAndPrintArtist(client, "1566") // Beyoncé
}

func requestAndPrintArtist(client *tidal.Client, artistID string) {
	fmt.Printf("Requesting artist with ID: %s\n", artistID)

	artist, err := client.GetArtist(context.Background(), artistID)
	if err != nil {
		log.Printf("Failed to get artist %s: %v", artistID, err)
		return
	}

	data, _ := json.MarshalIndent(artist, "", "  ")
	fmt.Printf("Artist Response for %s:\n%s\n", artistID, string(data))
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
