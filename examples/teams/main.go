package main

import (
	"context"
	"log"
	"strings"

	scalr "github.com/scalr/go-scalr"
)

func main() {
	accID := "acc-svrcncgh453bi8g"

	config := scalr.DefaultConfig()
	client, err := scalr.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	// List all teams in account
	tl, err := client.Teams.List(ctx, scalr.TeamListOptions{
		Account: scalr.String(accID),
	})
	if err != nil {
		log.Fatal(err)
	}

	if tl.TotalCount == 0 {
		log.Printf("No teams found in account %s", accID)
	} else {
		var teams []string
		for _, t := range tl.Items {
			teams = append(teams, t.Name)
		}
		log.Printf("Teams in account %s: %s", accID, strings.Join(teams, ", "))
	}

	// Create a new team
	t, err := client.Teams.Create(ctx, scalr.TeamCreateOptions{
		Name:        scalr.String("dev"),
		Description: scalr.String("Developers"),
		Account:     &scalr.Account{ID: accID},
		Users: []*scalr.User{
			{ID: "user-suh84u6vuvidtbg"},
			{ID: "user-suh84u72gfrbd30"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Update a team
	t, err = client.Teams.Update(ctx, t.ID, scalr.TeamUpdateOptions{
		Name: scalr.String("dev-new"),
		Users: []*scalr.User{
			{ID: "user-svrcmmpcrkmit1g"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
