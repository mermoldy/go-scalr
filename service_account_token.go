package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ ServiceAccountTokens = (*serviceAccountTokens)(nil)

// ServiceAccountTokens describes all the access token related methods that the
// Scalr IACP API supports.
type ServiceAccountTokens interface {
	// List service account's access tokens
	List(ctx context.Context, serviceAccountID string, options AccessTokenListOptions) (*AccessTokenList, error)
	// Create new access token for service account
	Create(ctx context.Context, serviceAccountID string, options AccessTokenCreateOptions) (*AccessToken, error)
}

// serviceAccountTokens implements ServiceAccountTokens.
type serviceAccountTokens struct {
	client *Client
}

// List the access tokens of ServiceAccount.
func (s *serviceAccountTokens) List(
	ctx context.Context, serviceAccountID string, options AccessTokenListOptions,
) (*AccessTokenList, error) {
	req, err := s.client.newRequest(
		"GET",
		fmt.Sprintf("service-accounts/%s/access-tokens", url.QueryEscape(serviceAccountID)),
		&options,
	)
	if err != nil {
		return nil, err
	}

	atl := &AccessTokenList{}
	err = s.client.do(ctx, req, atl)
	if err != nil {
		return nil, err
	}

	return atl, nil
}

// Create is used to create a new AccessToken for ServiceAccount.
func (s *serviceAccountTokens) Create(
	ctx context.Context, serviceAccountID string, options AccessTokenCreateOptions,
) (*AccessToken, error) {

	// Make sure we don't send a user provided ID.
	options.ID = ""

	if !validStringID(&serviceAccountID) {
		return nil, errors.New("invalid value for service account ID")
	}

	req, err := s.client.newRequest(
		"POST",
		fmt.Sprintf("service-accounts/%s/access-tokens", url.QueryEscape(serviceAccountID)),
		&options,
	)
	if err != nil {
		return nil, err
	}

	at := &AccessToken{}
	err = s.client.do(ctx, req, at)
	if err != nil {
		return nil, err
	}

	return at, nil
}
