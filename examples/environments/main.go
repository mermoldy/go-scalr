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

	// Create a new environment
	options := scalr.EnvironmentCreateOptions{
		Name:  scalr.String("example"),
		Email: scalr.String("info@example.com"),
	}

	org, err := client.Environments.Create(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	// Delete an environment
	err = client.Environments.Delete(ctx, org.Name)
	if err != nil {
		log.Fatal(err)
	}
}
