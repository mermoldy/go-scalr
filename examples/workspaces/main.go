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

	// Set your environment ID
	environmentID := "env-..."

	// Create a new workspace
	w, err := client.Workspaces.Create(ctx, scalr.WorkspaceCreateOptions{
		Name:        scalr.String("example-ws"),
		Environment: &scalr.Environment{ID: environmentID},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Update the workspace
	w, err = client.Workspaces.Update(ctx, w.ID, scalr.WorkspaceUpdateOptions{
		AutoApply:        scalr.Bool(false),
		TerraformVersion: scalr.String("0.12.28"),
		WorkingDirectory: scalr.String("my-app/infra"),
	})
	if err != nil {
		log.Fatal(err)
	}
}
