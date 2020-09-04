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
// Scalr API supports.
type Environments interface {
	// List all the environments visible to the current user.
	List(ctx context.Context, options EnvironmentListOptions) (*EnvironmentList, error)

	// Create a new environment with the given options.
	Create(ctx context.Context, options EnvironmentCreateOptions) (*Environment, error)

	// Read an environment by its name.
	Read(ctx context.Context, environment string) (*Environment, error)

	// Update attributes of an existing environment.
	Update(ctx context.Context, environment string, options EnvironmentUpdateOptions) (*Environment, error)

	// Delete an environment by its name.
	Delete(ctx context.Context, environment string) error

	// Capacity shows the current run capacity of an environment.
	Capacity(ctx context.Context, environment string) (*Capacity, error)

	// Entitlements shows the entitlements of an environment.
	Entitlements(ctx context.Context, environment string) (*Entitlements, error)

	// RunQueue shows the current run queue of an environment.
	RunQueue(ctx context.Context, environment string, options RunQueueOptions) (*RunQueue, error)
}

// environments implements Environments.
type environments struct {
	client *Client
}

// AuthPolicyType represents an authentication policy type.
type AuthPolicyType string

// List of available authentication policies.
const (
	AuthPolicyPassword  AuthPolicyType = "password"
	AuthPolicyTwoFactor AuthPolicyType = "two_factor_mandatory"
)

// EnterprisePlanType represents an enterprise plan type.
type EnterprisePlanType string

// List of available enterprise plan types.
const (
	EnterprisePlanDisabled EnterprisePlanType = "disabled"
	EnterprisePlanPremium  EnterprisePlanType = "premium"
	EnterprisePlanPro      EnterprisePlanType = "pro"
	EnterprisePlanTrial    EnterprisePlanType = "trial"
)

// EnvironmentList represents a list of environments.
type EnvironmentList struct {
	*Pagination
	Items []*Environment
}

// Environment represents a Scalr environment.
type Environment struct {
	Name                   string                  `jsonapi:"primary,organizations"`
	CollaboratorAuthPolicy AuthPolicyType          `jsonapi:"attr,collaborator-auth-policy"`
	CostEstimationEnabled  bool                    `jsonapi:"attr,cost-estimation-enabled"`
	CreatedAt              time.Time               `jsonapi:"attr,created-at,iso8601"`
	Email                  string                  `jsonapi:"attr,email"`
	EnterprisePlan         EnterprisePlanType      `jsonapi:"attr,enterprise-plan"`
	OwnersTeamSAMLRoleID   string                  `jsonapi:"attr,owners-team-saml-role-id"`
	Permissions            *EnvironmentPermissions `jsonapi:"attr,permissions"`
	SAMLEnabled            bool                    `jsonapi:"attr,saml-enabled"`
	SessionRemember        int                     `jsonapi:"attr,session-remember"`
	SessionTimeout         int                     `jsonapi:"attr,session-timeout"`
	TrialExpiresAt         time.Time               `jsonapi:"attr,trial-expires-at,iso8601"`
	TwoFactorConformant    bool                    `jsonapi:"attr,two-factor-conformant"`
}

// Capacity represents the current run capacity of an environment.
type Capacity struct {
	Environment string `jsonapi:"primary,organization-capacity"`
	Pending     int    `jsonapi:"attr,pending"`
	Running     int    `jsonapi:"attr,running"`
}

// Entitlements represents the entitlements of an environment.
type Entitlements struct {
	ID                    string `jsonapi:"primary,entitlement-sets"`
	Operations            bool   `jsonapi:"attr,operations"`
	PrivateModuleRegistry bool   `jsonapi:"attr,private-module-registry"`
	Sentinel              bool   `jsonapi:"attr,sentinel"`
	StateStorage          bool   `jsonapi:"attr,state-storage"`
	Teams                 bool   `jsonapi:"attr,teams"`
	VCSIntegrations       bool   `jsonapi:"attr,vcs-integrations"`
}

// RunQueue represents the current run queue of an environment.
type RunQueue struct {
	*Pagination
	Items []*Run
}

// EnvironmentPermissions represents the environment permissions.
type EnvironmentPermissions struct {
	CanCreateTeam               bool `json:"can-create-team"`
	CanCreateWorkspace          bool `json:"can-create-workspace"`
	CanCreateWorkspaceMigration bool `json:"can-create-workspace-migration"`
	CanDestroy                  bool `json:"can-destroy"`
	CanTraverse                 bool `json:"can-traverse"`
	CanUpdate                   bool `json:"can-update"`
	CanUpdateAPIToken           bool `json:"can-update-api-token"`
	CanUpdateOAuth              bool `json:"can-update-oauth"`
	CanUpdateSentinel           bool `json:"can-update-sentinel"`
}

// EnvironmentListOptions represents the options for listing environments.
type EnvironmentListOptions struct {
	ListOptions
}

// List all the environments visible to the current user.
func (s *environments) List(ctx context.Context, options EnvironmentListOptions) (*EnvironmentList, error) {
	req, err := s.client.newRequest("GET", "environments", &options)
	if err != nil {
		return nil, err
	}

	orgl := &EnvironmentList{}
	err = s.client.do(ctx, req, orgl)
	if err != nil {
		return nil, err
	}

	return orgl, nil
}

// EnvironmentCreateOptions represents the options for creating an environment.
type EnvironmentCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,organizations"`

	// Name of the environment.
	Name *string `jsonapi:"attr,name"`

	// Admin email address.
	Email *string `jsonapi:"attr,email"`

	// Session expiration (minutes).
	SessionRemember *int `jsonapi:"attr,session-remember,omitempty"`

	// Session timeout after inactivity (minutes).
	SessionTimeout *int `jsonapi:"attr,session-timeout,omitempty"`

	// Authentication policy.
	CollaboratorAuthPolicy *AuthPolicyType `jsonapi:"attr,collaborator-auth-policy,omitempty"`

	// Enable Cost Estimation
	CostEstimationEnabled *bool `jsonapi:"attr,cost-estimation-enabled,omitempty"`

	// The name of the "owners" team
	OwnersTeamSAMLRoleID *string `jsonapi:"attr,owners-team-saml-role-id,omitempty"`
}

func (o EnvironmentCreateOptions) valid() error {
	if !validString(o.Name) {
		return errors.New("name is required")
	}
	if !validStringID(o.Name) {
		return errors.New("invalid value for name")
	}
	if !validString(o.Email) {
		return errors.New("email is required")
	}
	return nil
}

// Create a new environment with the given options.
func (s *environments) Create(ctx context.Context, options EnvironmentCreateOptions) (*Environment, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "environments", &options)
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

// Read an environment by its name.
func (s *environments) Read(ctx context.Context, environment string) (*Environment, error) {
	if !validStringID(&environment) {
		return nil, errors.New("invalid value for environment")
	}

	u := fmt.Sprintf("environments/%s", url.QueryEscape(environment))
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

// EnvironmentUpdateOptions represents the options for updating an environment.
type EnvironmentUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,organizations"`

	// New name for the environment.
	Name *string `jsonapi:"attr,name,omitempty"`

	// New admin email address.
	Email *string `jsonapi:"attr,email,omitempty"`

	// Session expiration (minutes).
	SessionRemember *int `jsonapi:"attr,session-remember,omitempty"`

	// Session timeout after inactivity (minutes).
	SessionTimeout *int `jsonapi:"attr,session-timeout,omitempty"`

	// Authentication policy.
	CollaboratorAuthPolicy *AuthPolicyType `jsonapi:"attr,collaborator-auth-policy,omitempty"`

	// Enable Cost Estimation
	CostEstimationEnabled *bool `jsonapi:"attr,cost-estimation-enabled,omitempty"`

	// The name of the "owners" team
	OwnersTeamSAMLRoleID *string `jsonapi:"attr,owners-team-saml-role-id,omitempty"`
}

// Update attributes of an existing environment.
func (s *environments) Update(ctx context.Context, environment string, options EnvironmentUpdateOptions) (*Environment, error) {
	if !validStringID(&environment) {
		return nil, errors.New("invalid value for environment")
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("environments/%s", url.QueryEscape(environment))
	req, err := s.client.newRequest("PATCH", u, &options)
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

// Delete an environment by its name.
func (s *environments) Delete(ctx context.Context, environment string) error {
	if !validStringID(&environment) {
		return errors.New("invalid value for environment")
	}

	u := fmt.Sprintf("environments/%s", url.QueryEscape(environment))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Capacity shows the currently used capacity of an environment.
func (s *environments) Capacity(ctx context.Context, environment string) (*Capacity, error) {
	if !validStringID(&environment) {
		return nil, errors.New("invalid value for environment")
	}

	u := fmt.Sprintf("environments/%s/capacity", url.QueryEscape(environment))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	c := &Capacity{}
	err = s.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Entitlements shows the entitlements of an environment.
func (s *environments) Entitlements(ctx context.Context, environment string) (*Entitlements, error) {
	if !validStringID(&environment) {
		return nil, errors.New("invalid value for environment")
	}

	u := fmt.Sprintf("environments/%s/entitlement-set", url.QueryEscape(environment))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	e := &Entitlements{}
	err = s.client.do(ctx, req, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// RunQueueOptions represents the options for showing the queue.
type RunQueueOptions struct {
	ListOptions
}

// RunQueue shows the current run queue of an environment.
func (s *environments) RunQueue(ctx context.Context, environment string, options RunQueueOptions) (*RunQueue, error) {
	if !validStringID(&environment) {
		return nil, errors.New("invalid value for environment")
	}

	u := fmt.Sprintf("environments/%s/runs/queue", url.QueryEscape(environment))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	rq := &RunQueue{}
	err = s.client.do(ctx, req, rq)
	if err != nil {
		return nil, err
	}

	return rq, nil
}
