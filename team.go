package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Teams = (*teams)(nil)

// Teams describes all the team related methods that the
// Scalr API supports.
type Teams interface {
	List(ctx context.Context, options TeamListOptions) (*TeamList, error)
	Create(ctx context.Context, options TeamCreateOptions) (*Team, error)
	Read(ctx context.Context, teamID string) (*Team, error)
	Update(ctx context.Context, teamID string, options TeamUpdateOptions) (*Team, error)
	Delete(ctx context.Context, teamID string) error
}

// teams implements Teams.
type teams struct {
	client *Client
}

type Team struct {
	ID          string `jsonapi:"primary,teams"`
	Name        string `jsonapi:"attr,name,omitempty"`
	Description string `jsonapi:"attr,description,omitempty"`

	// Relations
	Account          *Account          `jsonapi:"relation,account"`
	IdentityProvider *IdentityProvider `jsonapi:"relation,identity-provider"`
	Users            []*User           `jsonapi:"relation,users"`
}

// TeamList represents a list of teams.
type TeamList struct {
	*Pagination
	Items []*Team
}

// TeamListOptions represents the options for listing teams.
type TeamListOptions struct {
	ListOptions

	Team             *string `url:"filter[team],omitempty"`
	Name             *string `url:"filter[name],omitempty"`
	Account          *string `url:"filter[account],omitempty"`
	IdentityProvider *string `url:"filter[identity-provider],omitempty"`
	Query            *string `url:"query,omitempty"`
	Sort             *string `url:"sort,omitempty"`
	Include          *string `url:"include,omitempty"`
}

// TeamCreateOptions represents the options for creating a new team.
type TeamCreateOptions struct {
	ID          string  `jsonapi:"primary,teams"`
	Name        *string `jsonapi:"attr,name"`
	Description *string `jsonapi:"attr,description"`

	// Relations
	Account          *Account          `jsonapi:"relation,account,omitempty"`
	IdentityProvider *IdentityProvider `jsonapi:"relation,identity-provider,omitempty"`
	Users            []*User           `jsonapi:"relation,users,omitempty"`
}

func (o TeamCreateOptions) valid() error {
	if !validString(o.Name) {
		return errors.New("name is required")
	}
	if o.Account != nil && !validStringID(&o.Account.ID) {
		return errors.New("invalid value for account ID")
	}
	if o.IdentityProvider != nil && !validStringID(&o.IdentityProvider.ID) {
		return errors.New("invalid value for identity provider ID")
	}

	return nil
}

// TeamUpdateOptions represents the options for updating a team.
type TeamUpdateOptions struct {
	ID          string  `jsonapi:"primary,teams"`
	Name        *string `jsonapi:"attr,name,omitempty"`
	Description *string `jsonapi:"attr,description,omitempty"`

	// Relations
	Users []*User `jsonapi:"relation,users"`
}

// List all the teams.
func (s *teams) List(ctx context.Context, options TeamListOptions) (*TeamList, error) {
	req, err := s.client.newRequest("GET", "teams", &options)
	if err != nil {
		return nil, err
	}

	tl := &TeamList{}
	err = s.client.do(ctx, req, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// Create a new team.
func (s *teams) Create(ctx context.Context, options TeamCreateOptions) (*Team, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	// Make sure we don't send a user provided ID.
	options.ID = ""
	req, err := s.client.newRequest("POST", "teams", &options)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = s.client.do(ctx, req, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Read team by its ID.
func (s *teams) Read(ctx context.Context, teamID string) (*Team, error) {
	if !validStringID(&teamID) {
		return nil, errors.New("invalid value for team ID")
	}

	u := fmt.Sprintf("teams/%s", url.QueryEscape(teamID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = s.client.do(ctx, req, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Update settings of an existing team.
func (s *teams) Update(ctx context.Context, teamID string, options TeamUpdateOptions) (*Team, error) {
	if !validStringID(&teamID) {
		return nil, errors.New("invalid value for team ID")
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("teams/%s", url.QueryEscape(teamID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	t := &Team{}
	err = s.client.do(ctx, req, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Delete team by its ID.
func (s *teams) Delete(ctx context.Context, teamID string) error {
	if !validStringID(&teamID) {
		return errors.New("invalid value for team ID")
	}

	u := fmt.Sprintf("teams/%s", url.QueryEscape(teamID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
