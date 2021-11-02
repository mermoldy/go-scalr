package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Users = (*users)(nil)

// Users describes all the user related methods that the
// Scalr API supports.
type Users interface {
	List(ctx context.Context, options UserListOptions) (*UserList, error)
	Read(ctx context.Context, userID string) (*User, error)
}

// users implements Users.
type users struct {
	client *Client
}

// UserStatus represents a user status.
type UserStatus string

// List of available user statuses.
const (
	UserStatusActive   UserStatus = "Active"
	UserStatusInactive UserStatus = "Inactive"
	UserStatusPending  UserStatus = "Pending"
)

// User represents a Scalr IAM user.
type User struct {
	ID          string     `jsonapi:"primary,users"`
	Status      UserStatus `jsonapi:"attr,status,omitempty"`
	Email       string     `jsonapi:"attr,email,omitempty"`
	Username    string     `jsonapi:"attr,username,omitempty"`
	FullName    string     `jsonapi:"attr,full-name,omitempty"`
	CreatedAt   time.Time  `jsonapi:"attr,created-at,iso8601"`
	LastLoginAt time.Time  `jsonapi:"attr,last-login-at,iso8601"`

	// Relations
	Teams             []*Team             `jsonapi:"relation,teams"`
	IdentityProviders []*IdentityProvider `jsonapi:"relation,identity-providers"`
}

// UserList represents a list of users.
type UserList struct {
	*Pagination
	Items []*User
}

// UserListOptions represents the options for listing users.
type UserListOptions struct {
	ListOptions

	User             *string `url:"filter[user],omitempty"`
	Email            *string `url:"filter[email],omitempty"`
	IdentityProvider *string `url:"filter[identity-provider],omitempty"`
	Query            *string `url:"query,omitempty"`
	Sort             *string `url:"sort,omitempty"`
	Include          *string `url:"include,omitempty"`
}

// List all the users.
func (s *users) List(ctx context.Context, options UserListOptions) (*UserList, error) {
	req, err := s.client.newRequest("GET", "users", &options)
	if err != nil {
		return nil, err
	}

	ul := &UserList{}
	err = s.client.do(ctx, req, ul)
	if err != nil {
		return nil, err
	}

	return ul, nil
}

// Read user by its ID.
func (s *users) Read(ctx context.Context, userID string) (*User, error) {
	if !validStringID(&userID) {
		return nil, errors.New("invalid value for user ID")
	}

	u := fmt.Sprintf("users/%s", url.QueryEscape(userID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	usr := &User{}
	err = s.client.do(ctx, req, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
