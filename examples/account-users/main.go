package main

import (
	"context"
	"log"
	"strings"

	scalr "github.com/scalr/go-scalr"
)

func main() {
	accID := "acc-svrcncgh453bi8g"
	userID := "user-suh84u6vhn64l0o"

	config := scalr.DefaultConfig()
	client, err := scalr.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	// List users that have access to an account
	aul, err := client.AccountUsers.List(ctx, scalr.AccountUserListOptions{
		Account: scalr.String(accID),
		Include: scalr.String("user"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(aul.Items) == 0 {
		log.Printf("No users found in account %s", accID)
	} else {
		var active []string
		for _, usr := range aul.Items {
			if usr.Status == scalr.AccountUserStatusActive {
				active = append(active, usr.User.Username)
			}
		}
		if len(active) == 0 {
			log.Printf("No active relations found for account %s", accID)
		} else {
			log.Printf("Active users in account %s: %s", accID, strings.Join(active, ", "))
		}
	}

	// List accounts the user has access to
	aul, err = client.AccountUsers.List(ctx, scalr.AccountUserListOptions{
		User:    scalr.String(userID),
		Include: scalr.String("account"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(aul.Items) == 0 {
		log.Printf("No accounts found for user %s", userID)
	} else {
		var active []string
		for _, usr := range aul.Items {
			if usr.Status == scalr.AccountUserStatusActive {
				active = append(active, usr.Account.Name)
			}
		}
		if len(active) == 0 {
			log.Printf("No active accounts found for user %s", userID)
		} else {
			log.Printf("Active accounts for user %s: %s", userID, strings.Join(active, ", "))
		}
	}
}
