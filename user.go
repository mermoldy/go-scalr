package scalr

import (
	"context"
)

// Compile-time proof of interface implementation.
var _ Users = (*users)(nil)

// Users describes all the user related methods that the Scalr API supports.
type Users interface {
	// ReadCurrent reads the details of the currently authenticated user.
	ReadCurrent(ctx context.Context) (*User, error)
}

// users implements Users.
type users struct {
	client *Client
}

// User represents a Scalr user.
type User struct {
	ID       string `jsonapi:"primary,users"`
	Email    string `jsonapi:"attr,email"`
	Username string `jsonapi:"attr,username"`
	FullName string `jsonapi:"attr,full-name"`
	// Relations
	// AuthenticationTokens *AuthenticationTokens `jsonapi:"relation,authentication-tokens"`
}

// ReadCurrent reads the details of the currently authenticated user.
func (s *users) ReadCurrent(ctx context.Context) (*User, error) {
	req, err := s.client.newRequest("GET", "account/details", nil)
	if err != nil {
		return nil, err
	}

	u := &User{}
	err = s.client.do(ctx, req, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
