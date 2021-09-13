package scalr

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ AccessTokens = (*accessTokens)(nil)

// AccessTokens describes all the access token related methods that the
// Scalr IACP API supports.
type AccessTokens interface {
	Update(ctx context.Context, accessTokenID string, options AccessTokenUpdateOptions) (*AccessToken, error)
	Delete(ctx context.Context, accessTokenID string) error
}

// accessTokens implements AccessTokens.
type accessTokens struct {
	client *Client
}

// AccessTokenList represents a list of access tokens.
type AccessTokenList struct {
	*Pagination
	Items []*AccessToken
}

// AccessToken represents a Scalr access token.
type AccessToken struct {
	ID          string    `jsonapi:"primary,access-tokens"`
	CreatedAt   time.Time `jsonapi:"attr,created-at,iso8601"`
	Description string    `jsonapi:"attr,description"`
	Token       string    `jsonapi:"attr,token"`
}

// AccessTokenUpdateOptions represents the options for updating an AccessToken.
type AccessTokenUpdateOptions struct {
	ID          string  `jsonapi:"primary,access-tokens"`
	Description *string `jsonapi:"attr,description"`
}

// Update is used to update an AccessToken.
func (s *accessTokens) Update(ctx context.Context, accessTokenID string, options AccessTokenUpdateOptions) (*AccessToken, error) {

	// Make sure we don't send a user provided ID.
	options.ID = ""

	if !validStringID(&accessTokenID) {
		return nil, fmt.Errorf("invalid value for access token ID: '%s'", accessTokenID)
	}

	if !validString(options.Description) {
		return nil, fmt.Errorf("invalid value for description: '%s'", *options.Description)
	}

	req, err := s.client.newRequest("PATCH", fmt.Sprintf("access-tokens/%s", url.QueryEscape(accessTokenID)), &options)
	if err != nil {
		return nil, err
	}

	accessToken := &AccessToken{}
	err = s.client.do(ctx, req, accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

// Delete an access token by its ID.
func (s *accessTokens) Delete(ctx context.Context, accessTokenID string) error {
	if !validStringID(&accessTokenID) {
		return fmt.Errorf("invalid value for access token ID: '%s'", accessTokenID)
	}

	t := fmt.Sprintf("access-tokens/%s", url.QueryEscape(accessTokenID))
	req, err := s.client.newRequest("DELETE", t, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
