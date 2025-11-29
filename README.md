# Tidal API Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/tmaffia/tidal.svg)](https://pkg.go.dev/pkg/github.com/tmaffia/tidal)
[![CI](https://github.com/tmaffia/tidal/actions/workflows/ci.yml/badge.svg)](https://github.com/tmaffia/tidal/actions/workflows/ci.yml)

A Go library for the Tidal API.

## Reference

- [Tidal OpenAPI Specification](https://tidal-music.github.io/tidal-api-reference/tidal-api-oas.json)

## Installation

```bash
go get github.com/tmaffia/tidal
```

## Usage

```go
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

 config := tidal.Config{
  ClientID:     clientID,
  ClientSecret: clientSecret,
 }

 // Create a client authenticated via Client Credentials flow
 client := tidal.NewClientCredentialsClient(context.Background(), config)

 // Fetch an Artist
 artist, err := client.Artists.Get(context.Background(), "1000", tidal.WithCountry("US"))
 if err != nil {
  log.Fatalf("Error fetching artist: %v", err)
 }

 fmt.Printf("Artist: %s (ID: %s)\n", artist.Name, artist.ID)
}
```

## Generating API Credentials

To use the Tidal API, you need to register an application on the Tidal Developer Portal.

1. Go to the [Tidal Developer Portal](https://developer.tidal.com/).
2. Log in with your Tidal account.
3. Navigate to the **Dashboard** or **My Apps**.
4. Click **Create App**.
5. Fill in the required details for your application.
6. Once created, you will see your **Client ID** and **Client Secret**.
7. Use these credentials to authenticate your client.

> **Note**: For the Client Credentials flow (used in the example above), ensure your app has the necessary permissions for the data you intend to access.

## Contributing

The linter runs automatically when a Pull Request is opened. Contributors are responsible for ensuring that all tests pass and linting succeeds.

To run the linter locally:

```bash
make lint
```

To run tests locally:

```bash
make test
```
