package main

import (
	"context"
	"log"
	"net/http"

	scalr "github.com/scalr/go-scalr"
)

func main() {
	config := &scalr.Config{
		Address:  "https://<example>.scalr.io",
		BasePath: "/api/iacp/v3/",
		Token:    "<your token>",
		Headers:  make(http.Header),
	}
	config.Headers.Set("Prefer", "profile=internal")

	client, err := scalr.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	// Get an environment
	env, err := client.Environments.Read(ctx, "env-...")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Environment created at %v", env.CreatedAt)
}
