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

For a quick start guide, please refer to the official documentation:
[Tidal API SDK Quick Start](https://developer.tidal.com/documentation/api-sdk/api-sdk-quick-start)

> **Important**: You must configure your **Redirect URIs** and **Scopes** in the Tidal Developer Portal for authentication to work correctly.
>
> - **Redirect URIs**: Ensure they exactly match what you use in your code (e.g., `http://localhost:8080/callback`).
> - **Scopes**: Select the scopes your application requires (e.g., `user.read`).

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
