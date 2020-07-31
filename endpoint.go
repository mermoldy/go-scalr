package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Endpoints = (*endpoints)(nil)

// Endpoints describes all the endpoints related methods that the Scalr
// IACP API supports.
//
// IACP API docs: TODO
type Endpoints interface {
	// List the endpoints.
	List(ctx context.Context, options EndpointListOptions) (*EndpointList, error)
	Create(ctx context.Context, options EndpointCreateOptions) (*Endpoint, error)
	Read(ctx context.Context, endpoint string) (*Endpoint, error)
	Update(ctx context.Context, endpoint string, options EndpointUpdateOptions) (*Endpoint, error)
	Delete(ctx context.Context, endpoint string) error
}

// endpoints implements Endpoints.
type endpoints struct {
	client *Client
}

// EndpointList represents a list of endpoints.
type EndpointList struct {
	*Pagination
	Items []*Endpoint
}

// Endpoint represents a Scalr IACP endpoint.
type Endpoint struct {
	ID          string `jsonapi:"primary,endpoints"`
	HTTPMethod  string `jsonapi:"attr,http-method"`
	MaxAttempts int    `jsonapi:"attr,max-attempts"`
	Name        string `jsonapi:"attr,name"`
	SecretKey   string `jsonapi:"attr,secret-key"`
	Timeout     int    `jsonapi:"attr,timeout"`

	// Relations
	Workspace   *Workspace   `jsonapi:"relation,workspace"`
	Environment *Environment `jsonapi:"relation,environment"`
	Account     *Account     `jsonapi:"relation,account"`
}

// EndpointListOptions represents the options for listing endpoints.
type EndpointListOptions struct {
	ListOptions

	// Query string.
	Query *string `url:"query,omitempty"`

	// The comma-separated list of attributes.
	Sort *string `url:"sort,omitempty"`

	// Scope filters.
	Workspace   *string `url:"filter[workspace],omitempty"`
	Environment *string `url:"filter[environment],omitempty"`
	Account     *string `url:"filter[account],omitempty"`
}

// List the endpoints.
func (s *endpoints) List(ctx context.Context, options EndpointListOptions) (*EndpointList, error) {
	u := fmt.Sprintf("endpoints")
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	wl := &EndpointList{}
	err = s.client.do(ctx, req, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

// EndpointCreateOptions represents the options for creating a new endpoint.
type EndpointCreateOptions struct {
	// For internal use only!
	ID          string  `jsonapi:"primary,endpoints"`
	HTTPMethod  *string `jsonapi:"attr,http-method"`
	MaxAttempts *int    `jsonapi:"attr,max-attempts"`
	Name        *string `jsonapi:"attr,name"`
	SecretKey   *string `jsonapi:"attr,secret-key"`
	Timeout     *int    `jsonapi:"attr,timeout"`

	// Relations
	Workspace   *Workspace   `jsonapi:"relation,workspace"`
	Environment *Environment `jsonapi:"relation,environment"`
	Account     *Account     `jsonapi:"relation,account"`
}

func (o EndpointCreateOptions) valid() error {
	if !validString(o.HTTPMethod) {
		return errors.New("HTTPMethod is required")
	}
	if !validStringID(&o.Workspace.ID) {
		return errors.New("invalid value for workspace ID")
	}
	return nil
}

// Create is used to create a new workspace.
func (s *endpoints) Create(ctx context.Context, options EndpointCreateOptions) (*Endpoint, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("endpoints")
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Endpoint{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Read a endpoint by its ID.
func (s *endpoints) Read(ctx context.Context, endpointID string) (*Endpoint, error) {
	if !validStringID(&endpointID) {
		return nil, errors.New("invalid value for endpoint ID")
	}

	u := fmt.Sprintf("endpoints/%s", url.QueryEscape(endpointID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Endpoint{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// EndpointUpdateOptions represents the options for updating an endpoint.
type EndpointUpdateOptions struct {
	ID          string  `jsonapi:"primary,endpoints"`
	HTTPMethod  *string `jsonapi:"attr,http-method"`
	MaxAttempts *int    `jsonapi:"attr,max-attempts"`
	SecretKey   *string `jsonapi:"attr,secret-key"`
	Timeout     *int    `jsonapi:"attr,timeout"`
}

// Update settings of an existing endpoint.
func (s *endpoints) Update(ctx context.Context, endpointID string, options EndpointUpdateOptions) (*Endpoint, error) {
	if !validStringID(&endpointID) {
		return nil, errors.New("invalid value for endpoint ID")
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("endpoints/%s", url.QueryEscape(endpointID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Endpoint{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Delete an endpoint by its ID.
func (s *endpoints) Delete(ctx context.Context, endpointID string) error {
	if !validStringID(&endpointID) {
		return errors.New("invalid value for endpoint ID")
	}

	u := fmt.Sprintf("endpoints/%s", url.QueryEscape(endpointID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
