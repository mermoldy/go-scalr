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

	// Create a new workspace
	w, err := client.Workspaces.Create(ctx, "env-name", scalr.WorkspaceCreateOptions{
		Name: scalr.String("my-app-tst"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Update the workspace
	w, err = client.Workspaces.Update(ctx, "env-name", w.Name, scalr.WorkspaceUpdateOptions{
		AutoApply:        scalr.Bool(false),
		TerraformVersion: scalr.String("0.12.0"),
		WorkingDirectory: scalr.String("my-app/infra"),
	})
	if err != nil {
		log.Fatal(err)
	}
}
