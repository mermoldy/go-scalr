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

	// Create a new policy group
	pg, err := client.PolicyGroups.Create(ctx, scalr.PolicyGroupCreateOptions{
		Name:        scalr.String("example-policy-group"),
		OpaVersion:  scalr.String("0.29.4"),
		Account:     &scalr.Account{ID: "acc-..."},
		VcsProvider: &scalr.VcsProvider{ID: "vcs-..."},
		VCSRepo: &scalr.PolicyGroupVCSRepoOptions{
			Identifier: scalr.String("foo/bar"),
			Branch:     scalr.String("dev"),
			Path:       scalr.String("opa"),
		},
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("The policy group with name %v was created. ID: %v", pg.Name, pg.ID)
	}

	// Update the policy group
	pg, err = client.PolicyGroups.Update(ctx, pg.ID, scalr.PolicyGroupUpdateOptions{
		Name: scalr.String("go-scalr-policy-group"),
		VCSRepo: &scalr.PolicyGroupVCSRepoOptions{
			Identifier: scalr.String("foo/bar"),
			Branch:     scalr.String("main"),
			Path:       scalr.String("opa"),
		},
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("The policy group with id %v was updated. New repo config: %+v", pg.ID, pg.VCSRepo)
	}

	// Delete the policy group
	err = client.PolicyGroups.Delete(ctx, pg.ID)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("The policy group with id %v was deleted.", pg.ID)
	}
}
