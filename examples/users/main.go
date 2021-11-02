package main

import (
	"context"
	"log"

	scalr "github.com/scalr/go-scalr"
)

func main() {
	config := scalr.DefaultConfig()
	client, err := scalr.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	// List all users
	ul, err := client.Users.List(ctx, scalr.UserListOptions{
		ListOptions: scalr.ListOptions{PageSize: 99},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Obtained %d users", ul.TotalCount)
}
