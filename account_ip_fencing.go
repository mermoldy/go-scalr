package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AccountIPAllowlists = (*accountIPAllowlists)(nil)

// AccountIPAllowlists describes methods for updating and reading ip fencing rules that the
// Scalr IACP API supports.
type AccountIPAllowlists interface {
	Read(ctx context.Context, account string) (*AccountIPAllowlist, error)
	Update(ctx context.Context, account string, options AccountIPAllowlistUpdateOptions) (*AccountIPAllowlist, error)
}

// accountIPAllowlists implements AccountIPAllowlists.
type accountIPAllowlists struct {
	client *Client
}

type AccountIPAllowlist struct {
	Account
	IPAllowlist []string `jsonapi:"attr,ip-allowlist"`
}

// Read a account by its ID.
func (s *accountIPAllowlists) Read(ctx context.Context, accountID string) (*AccountIPAllowlist, error) {
	if !validStringID(&accountID) {
		return nil, errors.New("invalid value for account ID")
	}

	u := fmt.Sprintf("accounts/%s", url.QueryEscape(accountID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	a := &AccountIPAllowlist{}
	err = s.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type AccountIPAllowlistUpdateOptions struct {
	IPAllowlist *[]string `json:"ip-allowlist,omitempty"`
}

func (s *accountIPAllowlists) Update(ctx context.Context, accountID string, options AccountIPAllowlistUpdateOptions) (*AccountIPAllowlist, error) {
	if !validStringID(&accountID) {
		return nil, errors.New("invalid value for account ID")
	}

	for _, network := range *options.IPAllowlist {
		if !validIPv4Network(&network) {
			return nil, fmt.Errorf("invalid value for ip allowlist entry: %s", network)
		}
	}

	u := fmt.Sprintf("accounts/%s/actions/set-ip-allowlist", url.QueryEscape(accountID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	a := &AccountIPAllowlist{}
	err = s.client.do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
