package scalr

import (
	"context"
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

// Environment represents a Terraform Enterprise environment.
type Environment struct {
	ID                    string            `jsonapi:"primary,environments"`
	Name                  string            `jsonapi:"attr,name"`
	CostEstimationEnabled bool              `jsonapi:"attr,cost-estimation-enabled"`
	CreatedAt             time.Time         `jsonapi:"attr,created-at,iso8601"`
	CreatedBy             string            `jsonapi:"attr,created-by"`
	Status                EnvironmentStatus `jsonapi:"attr,status"`

	// Relations
	Account *Account `jsonapi:"relation,account"`
}

// duplicate Environment to change type to "organizations" (workspace.Organization)
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

// Read an environment by its ID.
func (s *environments) Read(ctx context.Context, environmentID string) (*Environment, error) {
	if !validStringID(&environmentID) {
		return nil, fmt.Errorf("invalid value for environment ID: %v", environmentID)
	}

	u := fmt.Sprintf("environments/%s", url.QueryEscape(environmentID))
	req, err := s.client.newRequest("GET", u, nil)
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
