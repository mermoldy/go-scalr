package main

import (
	"context"
	"log"

	scalr "github.com/scalr/go-scalr"
)

func main() {
	config := &scalr.Config{
		Token: "insert-your-token-here",
	}

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
