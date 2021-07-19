package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AccessPolicies = (*accessPolicies)(nil)

// AccessPolicies describes all the accessPolicy related methods that the
// Scalr IACP API supports.
type AccessPolicies interface {
	List(ctx context.Context, options AccessPolicyListOptions) (*AccessPolicyList, error)
	Read(ctx context.Context, accessPolicyID string) (*AccessPolicy, error)
	Create(ctx context.Context, options AccessPolicyCreateOptions) (*AccessPolicy, error)
	Update(ctx context.Context, accessPolicyID string, options AccessPolicyUpdateOptions) (*AccessPolicy, error)
	Delete(ctx context.Context, accessPolicyID string) error
}

// accessPolicies implements AccessPolicies.
type accessPolicies struct {
	client *Client
}

// AccessPolicyList represents a list of accessPolicies.
type AccessPolicyList struct {
	*Pagination
	Items []*AccessPolicy
}

// AccessPolicy represents a Scalr accessPolicy.
type AccessPolicy struct {
	ID             string          `jsonapi:"primary,access-policies"`
	IsSystem       bool            `jsonapi:"attr,is-system"`
	Roles          []*Role         `jsonapi:"relation,roles"`
	User           *User           `jsonapi:"relation,user,omitempty"`
	Team           *Team           `jsonapi:"relation,team,omitempty"`
	ServiceAccount *ServiceAccount `jsonapi:"relation,service-account,omitempty"`
	Account        *Account        `jsonapi:"relation,account,omitempty"`
	Environment    *Environment    `jsonapi:"relation,environment,omitempty"`
	Workspace      *Workspace      `jsonapi:"relation,workspace,omitempty"`
}

// AccessPolicyCreateOptions represents the options for creating a new AccessPolicy.
type AccessPolicyCreateOptions struct {
	ID string `jsonapi:"primary,access-policies"`

	// Relations
	Roles []*Role `jsonapi:"relation,roles"`
	// The object of access policy, one of this fields must be filled
	User           *User           `jsonapi:"relation,user,omitempty"`
	Team           *Team           `jsonapi:"relation,team,omitempty"`
	ServiceAccount *ServiceAccount `jsonapi:"relation,service-account,omitempty"`
	// Scope
	Account     *Account     `jsonapi:"relation,account,omitempty"`
	Environment *Environment `jsonapi:"relation,environment,omitempty"`
	Workspace   *Workspace   `jsonapi:"relation,workspace,omitempty"`
}

func (o AccessPolicyCreateOptions) valid() error {
	if len(o.Roles) == 0 {
		return errors.New("at least one role must be provided")
	}

	if o.Account == nil && o.Environment == nil && o.Workspace == nil {
		return errors.New("one of: account,environment,workspace must be provided")
	}

	var scopeId, field string
	if o.Account != nil {
		scopeId = o.Account.ID
		field = "account"
	} else if o.Environment != nil {
		scopeId = o.Environment.ID
		field = "environment"
	} else {
		scopeId = o.Workspace.ID
		field = "workspace"
	}

	if !validStringID(&scopeId) {
		return fmt.Errorf("invalid value for %v ID: %v", field, scopeId)
	}

	if o.User == nil && o.Team == nil && o.ServiceAccount == nil {
		return errors.New("one of: user,team,service_account must be provided")
	}

	var object string
	if o.User != nil {
		object = o.User.ID
		field = "user"
	} else if o.Team != nil {
		object = o.Team.ID
		field = "team"
	} else {
		object = o.ServiceAccount.ID
		field = "service_account"
	}

	if !validStringID(&object) {
		return fmt.Errorf("invalid value for %v ID: %v", field, object)
	}

	return nil
}

// AccessPolicyListOptions represents the options for listing access policies.
type AccessPolicyListOptions struct {
	ListOptions

	Environment    *string `url:"filter[environment],omitempty"`
	Account        *string `url:"filter[account],omitempty"`
	Workspace      *string `url:"filter[workspace],omitempty"`
	User           *string `url:"filter[user],omitempty"`
	ServiceAccount *string `url:"filter[service-account],omitempty"`
	Team           *string `url:"filter[team],omitempty"`
	Include        string  `url:"include,omitempty"`
}

// List the accessPolicies.
func (s *accessPolicies) List(ctx context.Context, options AccessPolicyListOptions) (*AccessPolicyList, error) {
	req, err := s.client.newRequest("GET", "access-policies", &options)
	if err != nil {
		return nil, err
	}

	accessPolicyl := &AccessPolicyList{}
	err = s.client.do(ctx, req, accessPolicyl)
	if err != nil {
		return nil, err
	}

	return accessPolicyl, nil
}

// Create is used to create a new AccessPolicy.
func (s *accessPolicies) Create(ctx context.Context, options AccessPolicyCreateOptions) (*AccessPolicy, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	// Make sure we don't send a user provided ID.
	options.ID = ""
	req, err := s.client.newRequest("POST", "access-policies", &options)
	if err != nil {
		return nil, err
	}

	accessPolicy := &AccessPolicy{}
	err = s.client.do(ctx, req, accessPolicy)
	if err != nil {
		return nil, err
	}

	return accessPolicy, nil
}

// Read an accessPolicy by its ID.
func (s *accessPolicies) Read(ctx context.Context, accessPolicyID string) (*AccessPolicy, error) {
	if !validStringID(&accessPolicyID) {
		return nil, errors.New("invalid value for accessPolicy")
	}

	u := fmt.Sprintf("access-policies/%s", url.QueryEscape(accessPolicyID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	accessPolicy := &AccessPolicy{}
	err = s.client.do(ctx, req, accessPolicy)
	if err != nil {
		return nil, err
	}

	return accessPolicy, nil
}

// AccessPolicyUpdateOptions represents the options for updating an accessPolicy.
type AccessPolicyUpdateOptions struct {
	// For internal use only!
	ID    string  `jsonapi:"primary,access-policies"`
	Roles []*Role `jsonapi:"relation,roles"`
}

// Update settings of an existing accessPolicy.
func (s *accessPolicies) Update(ctx context.Context, accessPolicyID string, options AccessPolicyUpdateOptions) (*AccessPolicy, error) {
	// Make sure we don't send a user provided ID.
	options.ID = ""

	if len(options.Roles) == 0 {
		return nil, errors.New("at least one role must be provided")
	}

	u := fmt.Sprintf("access-policies/%s", url.QueryEscape(accessPolicyID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	accessPolicy := &AccessPolicy{}
	err = s.client.do(ctx, req, accessPolicy)
	if err != nil {
		return nil, err
	}

	return accessPolicy, nil
}

// Delete an accessPolicy by its ID.
func (s *accessPolicies) Delete(ctx context.Context, accessPolicyID string) error {
	if !validStringID(&accessPolicyID) {
		return errors.New("invalid value for accessPolicy ID")
	}

	u := fmt.Sprintf("access-policies/%s", url.QueryEscape(accessPolicyID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
