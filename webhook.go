package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ Webhooks = (*webhooks)(nil)

// Webhooks describes all the webhooks related methods that the Scalr
// IACP API supports.
//
// IACP API docs: https://www.scalr.com/docs/en/latest/api/index.html
type Webhooks interface {
	// List the webhooks.
	List(ctx context.Context, options WebhookListOptions) (*WebhookList, error)
	Create(ctx context.Context, options WebhookCreateOptions) (*Webhook, error)
	Read(ctx context.Context, webhook string) (*Webhook, error)
	Update(ctx context.Context, webhook string, options WebhookUpdateOptions) (*Webhook, error)
	Delete(ctx context.Context, webhook string) error
}

// webhooks implements Webhooks.
type webhooks struct {
	client *Client
}

// WebhookList represents a list of webhooks.
type WebhookList struct {
	*Pagination
	Items []*Webhook
}

type EventDefinition struct {
	ID string `jsonapi:"primary,event-definitions"`
}

// Webhook represents a Scalr IACP webhook.
type Webhook struct {
	ID              string     `jsonapi:"primary,webhooks"`
	Enabled         bool       `jsonapi:"attr,enabled"`
	LastTriggeredAt *time.Time `jsonapi:"attr,last-triggered-at,iso8601"`
	Name            string     `jsonapi:"attr,name"`

	// Relations
	Workspace   *Workspace         `jsonapi:"relation,workspace"`
	Environment *Environment       `jsonapi:"relation,environment"`
	Account     *Account           `jsonapi:"relation,account"`
	Endpoint    *Endpoint          `jsonapi:"relation,endpoint"`
	Events      []*EventDefinition `jsonapi:"relation,events"`
}

// WebhookListOptions represents the options for listing webhooks.
type WebhookListOptions struct {
	ListOptions

	// Query string.
	Query *string `url:"query,omitempty"`

	// The comma-separated list of attributes.
	Sort *string `url:"sort,omitempty"`

	// The comma-separated list of relationship paths.
	Include *string `url:"include,omitempty"`

	// Filter by enabled
	Enabled *bool `url:"filter[webhook][enabled],omitempty"`

	// Event filter
	Event *string `url:"filter[event],omitempty"`

	// Scope filters.
	Workspace   *string `url:"filter[workspace],omitempty"`
	Environment *string `url:"filter[environment],omitempty"`
	Account     *string `url:"filter[account],omitempty"`
}

// List the webhooks.
func (s *webhooks) List(ctx context.Context, options WebhookListOptions) (*WebhookList, error) {
	req, err := s.client.newRequest("GET", "webhooks", &options)
	if err != nil {
		return nil, err
	}

	wl := &WebhookList{}
	err = s.client.do(ctx, req, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

// WebhookCreateOptions represents the options for creating a new webhook.
type WebhookCreateOptions struct {
	ID      string  `jsonapi:"primary,webhooks"`
	Enabled *bool   `jsonapi:"attr,enabled,omitempty"`
	Name    *string `jsonapi:"attr,name"`

	// Relations
	Workspace   *Workspace         `jsonapi:"relation,workspace,omitempty"`
	Environment *Environment       `jsonapi:"relation,environment,omitempty"`
	Account     *Account           `jsonapi:"relation,account"`
	Endpoint    *Endpoint          `jsonapi:"relation,endpoint"`
	Events      []*EventDefinition `jsonapi:"relation,events"`
}

func (o WebhookCreateOptions) valid() error {
	if o.Name == nil {
		return errors.New("missing name")
	}
	return nil
}

// Create is used to create a new webhook.
func (s *webhooks) Create(ctx context.Context, options WebhookCreateOptions) (*Webhook, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "webhooks", &options)
	if err != nil {
		return nil, err
	}

	w := &Webhook{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Read a webhook by its ID.
func (s *webhooks) Read(ctx context.Context, webhookID string) (*Webhook, error) {
	if !validStringID(&webhookID) {
		return nil, errors.New("invalid value for webhook ID")
	}

	u := fmt.Sprintf("webhooks/%s", url.QueryEscape(webhookID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Webhook{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// WebhookUpdateOptions represents the options for updating an webhook.
type WebhookUpdateOptions struct {
	// For internal use only!
	ID      string  `jsonapi:"primary,webhooks"`
	Enabled *bool   `jsonapi:"attr,enabled,omitempty"`
	Name    *string `jsonapi:"attr,name"`

	// Relations
	Endpoint *Endpoint          `jsonapi:"relation,endpoint"`
	Events   []*EventDefinition `jsonapi:"relation,events"`
}

// Update settings of an existing webhook.
func (s *webhooks) Update(ctx context.Context, webhookID string, options WebhookUpdateOptions) (*Webhook, error) {
	if !validStringID(&webhookID) {
		return nil, errors.New("invalid value for webhook ID")
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("webhooks/%s", url.QueryEscape(webhookID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Webhook{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Delete an webhook by its ID.
func (s *webhooks) Delete(ctx context.Context, webhookID string) error {
	if !validStringID(&webhookID) {
		return errors.New("invalid value for webhook ID")
	}

	u := fmt.Sprintf("webhooks/%s", url.QueryEscape(webhookID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
