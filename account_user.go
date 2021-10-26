package scalr

import (
	"context"
	"errors"
)

// Compile-time proof of interface implementation.
var _ AccountUsers = (*accountUsers)(nil)

// AccountUsers describes all the account user related methods that the
// Scalr IACP API supports.
type AccountUsers interface {
	List(ctx context.Context, options AccountUserListOptions) (*AccountUserList, error)
}

// accountUsers implements AccountUsers.
type accountUsers struct {
	client *Client
}

// AccountUserStatus represents a status of account user relation.
type AccountUserStatus string

// List of available policy group statuses.
const (
	AccountUserStatusActive   AccountUserStatus = "Active"
	AccountUserStatusInactive AccountUserStatus = "Inactive"
	AccountUserStatusPending  AccountUserStatus = "Pending"
)

// AccountUserListOptions represents the options for listing account users.
type AccountUserListOptions struct {
	Account *string `url:"filter[account],omitempty"`
	User    *string `url:"filter[user],omitempty"`
	Query   *string `url:"query,omitempty"`
	Sort    *string `url:"sort,omitempty"`
	Include *string `url:"include,omitempty"`
}

func (o AccountUserListOptions) validate() error {
	if !(validString(o.Account) || validString(o.User)) {
		return errors.New("either filter[account] or filter[user] is required")
	}
	return nil
}

// AccountUserList represents a list of account users.
type AccountUserList struct {
	*Pagination
	Items []*AccountUser
}

// AccountUser represents a Scalr account user.
type AccountUser struct {
	ID     string            `jsonapi:"primary,account-users"`
	Status AccountUserStatus `jsonapi:"attr,status"`

	// Relations
	Account *Account `jsonapi:"relation,account"`
	User    *User    `jsonapi:"relation,user"`
	Teams   []*Team  `jsonapi:"relation,teams"`
}

// List all the account users.
func (s *accountUsers) List(ctx context.Context, options AccountUserListOptions) (*AccountUserList, error) {
	if err := options.validate(); err != nil {
		return nil, err
	}

	req, err := s.client.newRequest("GET", "account-users", &options)
	if err != nil {
		return nil, err
	}

	aul := &AccountUserList{}
	err = s.client.do(ctx, req, aul)
	if err != nil {
		return nil, err
	}

	return aul, nil
}
