package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Environments = (*environments)(nil)

// Environments describes all the environment related methods that the
// Scalr IACP API supports.
type Environments interface {
	// Read an environment by its ID.
	Read(ctx context.Context, environmentID string) (*Environment, error)
	Create(ctx context.Context, options EnvironmentCreateOptions) (*Environment, error)
	Delete(ctx context.Context, environmentID string) error
}

// environments implements Environments.
type environments struct {
	client *Client
}

// EnvironmentStatus represents an environment status.
type EnvironmentStatus string

// List of available environment statuses.
const (
	EnvironmentStatusActive   EnvironmentStatus = "Active"
	EnvironmentStatusInactive EnvironmentStatus = "Inactive"
)

// CloudCredential relationship
type CloudCredential struct {
	ID string `jsonapi:"primary,cloud-credentials"`
}

// PolicyGroup relationship
type PolicyGroup struct {
	ID string `jsonapi:"primary,policy-groups"`
}

// Environment represents a Scalr environment.
type Environment struct {
	ID                    string            `jsonapi:"primary,environments"`
	Name                  string            `jsonapi:"attr,name"`
	CostEstimationEnabled bool              `jsonapi:"attr,cost-estimation-enabled"`
	CreatedAt             time.Time         `jsonapi:"attr,created-at,iso8601"`
	Status                EnvironmentStatus `jsonapi:"attr,status"`

	// Relations
	Account          *Account           `jsonapi:"relation,account"`
	CloudCredentials []*CloudCredential `jsonapi:"relation,cloud-credentials"`
	PolicyGroups     []*PolicyGroup     `jsonapi:"relation,policy-groups"`
	CreatedBy        *User              `jsonapi:"relation,created-by"`
}

// Organization is Environment included in Workspace - always prefer Environment
type Organization struct {
	ID                    string            `jsonapi:"primary,organizations"`
	Name                  string            `jsonapi:"attr,name"`
	CostEstimationEnabled bool              `jsonapi:"attr,cost-estimation-enabled"`
	CreatedAt             time.Time         `jsonapi:"attr,created-at,iso8601"`
	CreatedBy             string            `jsonapi:"attr,created-by"`
	Status                EnvironmentStatus `jsonapi:"attr,status"`

	// Relations
	Account *Account `jsonapi:"relation,account"`
}

// EnvironmentCreateOptions represents the options for creating a new Environment.
type EnvironmentCreateOptions struct {
	ID                    string  `jsonapi:"primary,environments"`
	Name                  *string `jsonapi:"attr,name"`
	CostEstimationEnabled *bool   `jsonapi:"attr,cost-estimation-enabled"`

	// Relations
	Account          *Account           `jsonapi:"relation,account"`
	CloudCredentials []*CloudCredential `jsonapi:"relation,cloud-credentials"`
	PolicyGroups     []*PolicyGroup     `jsonapi:"relation,policy-groups"`
}

// Create is used to create a new Environment.
func (s *environments) Create(ctx context.Context, options EnvironmentCreateOptions) (*Environment, error) {
	if !validStringID(&options.Account.ID) {
		return nil, errors.New("invalid value for Account.ID")
	}
	// if err := options.valid(); err != nil {
	// 	return nil, err
	// }

	// Make sure we don't send a user provided ID.
	options.ID = ""
	req, err := s.client.newRequest("POST", "environments", &options)
	if err != nil {
		return nil, err
	}

	environment := &Environment{}
	err = s.client.do(ctx, req, environment)
	if err != nil {
		return nil, err
	}

	return environment, nil
}

// Read an environment by its ID.
func (s *environments) Read(ctx context.Context, environmentID string) (*Environment, error) {
	if !validStringID(&environmentID) {
		return nil, fmt.Errorf("invalid value for environment ID: %v", environmentID)
	}

	options := struct {
		Include string `url:"include"`
	}{
		Include: "created-by",
	}
	u := fmt.Sprintf("environments/%s", url.QueryEscape(environmentID))
	req, err := s.client.newRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	org := &Environment{}
	err = s.client.do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Delete an environment by its ID.
func (s *environments) Delete(ctx context.Context, environmentID string) error {
	if !validStringID(&environmentID) {
		return errors.New("invalid value for environment ID")
	}

	u := fmt.Sprintf("environments/%s", url.QueryEscape(environmentID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
