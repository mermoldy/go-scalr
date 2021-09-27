package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ VcsProviders = (*vcs_providers)(nil)

// VcsProviders describes all the VcsProviders related methods that the Scalr
// IACP API supports.
//
// IACP API docs: https://www.scalr.com/docs/en/latest/api/index.html
type VcsProviders interface {
	// List the vcs providers.
	List(ctx context.Context, options VcsProvidersListOptions) (*VcsProvidersList, error)
	Create(ctx context.Context, options VcsProviderCreateOptions) (*VcsProvider, error)
	Read(ctx context.Context, vcsProvider string) (*VcsProvider, error)
	Update(ctx context.Context, vcsProvider string, options VcsProviderUpdateOptions) (*VcsProvider, error)
	Delete(ctx context.Context, vcsProvider string) error
}

// vcs_providers implements VcsProviders.
type vcs_providers struct {
	client *Client
}

// VcsType represents a type VCS provider.
type VcsType string

const (
	Github              VcsType = "github"
	GithubEnterprise    VcsType = "github_enterprise"
	Gitlab              VcsType = "gitlab"
	GitlabEnterprise    VcsType = "gitlab_enterprise"
	Bitbucket           VcsType = "bitbucket"
	BitbucketEnterprise VcsType = "bitbucket_enterprise"
	AzureDevOpsServices VcsType = "azure_dev_ops_services"
)

// AuthType represents the authorization type used in VCS provider.
type AuthType string

const (
	Oauth2        AuthType = "oauth2"
	PersonalToken AuthType = "personal_token"
)

// VcsProvidersList represents a list of VCS providers.
type VcsProvidersList struct {
	*Pagination
	Items []*VcsProvider
}

// OAuth contains the custom hooks field.
type OAuth struct {
	ClientId     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
}

// VcsProvider represents a Scalr IACP VcsProvider.
type VcsProvider struct {
	ID       string   `jsonapi:"primary,vcs-providers"`
	Name     string   `jsonapi:"attr,name"`
	Url      string   `jsonapi:"attr,url"`
	VcsType  VcsType  `jsonapi:"attr,vcs-type"`
	AuthType AuthType `jsonapi:"attr,auth-type"`
	OAuth    *OAuth   `jsonapi:"attr,oauth"`

	// Relations
	Environments []*Environment `jsonapi:"relation,environments"`
	Account      *Account       `jsonapi:"relation,account"`
}

// VcsProvidersListOptions represents the options for listing vcs providers.
type VcsProvidersListOptions struct {
	ListOptions

	// Query string.
	Query *string `url:"query,omitempty"`

	// The comma-separated list of attributes.
	Sort *string `url:"sort,omitempty"`

	// Filter by vcs-type
	VcsType *VcsType `url:"filter[vcs-type],omitempty"`

	// Scope filters.
	Environment *string `url:"filter[environment],omitempty"`
	Account     *string `url:"filter[account],omitempty"`
}

// List the vcs providers.
func (s *vcs_providers) List(ctx context.Context, options VcsProvidersListOptions) (*VcsProvidersList, error) {
	req, err := s.client.newRequest("GET", "vcs-providers", &options)
	if err != nil {
		return nil, err
	}

	providersList := &VcsProvidersList{}
	err = s.client.do(ctx, req, providersList)
	if err != nil {
		return nil, err
	}

	return providersList, nil
}

// VcsProviderCreateOptions represents the options for creating a new vcs provider.
type VcsProviderCreateOptions struct {
	ID       string   `jsonapi:"primary,vcs-providers"`
	Name     *string  `jsonapi:"attr,name"`
	VcsType  VcsType  `jsonapi:"attr,vcs-type"`
	AuthType AuthType `jsonapi:"attr,auth-type"`
	OAuth    *OAuth   `jsonapi:"attr,oauth"`
	Token    string   `jsonapi:"attr,token"`
	Url      *string  `jsonapi:"attr,url"`

	// Relations
	Environments []*Environment `jsonapi:"relation,environments,omitempty"`
	Account      *Account       `jsonapi:"relation,account,omitempty"`
}

func (o VcsProviderCreateOptions) valid() error {
	if o.Name == nil {
		return errors.New("missing name")
	}
	return nil
}

// Create is used to create a new vcs provider.
func (s *vcs_providers) Create(ctx context.Context, options VcsProviderCreateOptions) (*VcsProvider, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "vcs-providers", &options)
	if err != nil {
		return nil, err
	}

	w := &VcsProvider{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Read a vcs provider by its ID.
func (s *vcs_providers) Read(ctx context.Context, vcsProviderID string) (*VcsProvider, error) {
	if !validStringID(&vcsProviderID) {
		return nil, errors.New("invalid value for vcs provider ID")
	}

	u := fmt.Sprintf("vcs-providers/%s", url.QueryEscape(vcsProviderID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	w := &VcsProvider{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// VcsProviderUpdateOptions represents the options for updating a vcs provider.
type VcsProviderUpdateOptions struct {
	// For internal use only!
	Name  *string `jsonapi:"attr,name"`
	Token *string `jsonapi:"attr,token"`
	Url   *string `jsonapi:"attr,url"`
}

// Update settings of an existing vcs provider.
func (s *vcs_providers) Update(ctx context.Context, vcsProviderId string, options VcsProviderUpdateOptions) (*VcsProvider, error) {
	if !validStringID(&vcsProviderId) {
		return nil, errors.New("invalid value for vcs provider ID")
	}

	u := fmt.Sprintf("vcs-providers/%s", url.QueryEscape(vcsProviderId))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &VcsProvider{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Delete an vcs provider by its ID.
func (s *vcs_providers) Delete(ctx context.Context, vcsProviderId string) error {
	if !validStringID(&vcsProviderId) {
		return errors.New("invalid value for vcs provider ID")
	}

	u := fmt.Sprintf("vcs-providers/%s", url.QueryEscape(vcsProviderId))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
