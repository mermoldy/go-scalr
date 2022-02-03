package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Accounts = (*accounts)(nil)

// AccountIPAllowlists describes methods for updating and reading ip fencing rules that the
// Scalr IACP API supports.
type Accounts interface {
	Read(ctx context.Context, account string) (*Account, error)
	Update(ctx context.Context, account string, options AccountUpdateOptions) (*Account, error)
}

// accountIPAllowlists implements AccountIPAllowlists.
type accounts struct {
	client *Client
}

// Account represents a Scalr IACP account.
type Account struct {
	ID         string   `jsonapi:"primary,accounts"`
	Name       string   `jsonapi:"attr,name"`
	AllowedIPs []string `jsonapi:"attr,allowed-ips"`
}

// Read a account by its ID.
func (s *accounts) Read(ctx context.Context, accountID string) (*Account, error) {
	if !validStringID(&accountID) {
		return nil, errors.New("invalid value for account ID")
	}

	u := fmt.Sprintf("accounts/%s", url.QueryEscape(accountID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	a := &Account{}
	err = s.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type AccountUpdateOptions struct {
	ID         string    `jsonapi:"primary,accounts"`
	AllowedIPs *[]string `jsonapi:"attr,allowed-ips,omitempty"`
}

func (s *accounts) Update(ctx context.Context, accountID string, options AccountUpdateOptions) (*Account, error) {
	if !validStringID(&accountID) {
		return nil, errors.New("invalid value for account ID")
	}

	for _, network := range *options.AllowedIPs {
		if !validIPv4Network(&network) {
			return nil, fmt.Errorf("invalid value for ip allowlist entry: %s", network)
		}
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("accounts/%s", url.QueryEscape(accountID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	a := &Account{}
	err = s.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
