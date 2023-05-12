package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ SlackIntegrations = (*slackIntegrations)(nil)

// SlackIntegrations describes all the SlackIntegration related methods that the Scalr
// IACP API supports.
//
// IACP API docs: https://www.scalr.com/docs/en/latest/api/index.html
type SlackIntegrations interface {
	List(ctx context.Context, options SlackIntegrationListOptions) (*SlackIntegrationList, error)
	Create(ctx context.Context, options SlackIntegrationCreateOptions) (*SlackIntegration, error)
	Read(ctx context.Context, slackIntegration string) (*SlackIntegration, error)
	Update(ctx context.Context, slackIntegration string, options SlackIntegrationUpdateOptions) (*SlackIntegration, error)
	Delete(ctx context.Context, slackIntegration string) error
	GetConnection(ctx context.Context, accID string) (*SlackConnection, error)
	GetChannels(ctx context.Context, accID string, options SlackChannelListOptions) (*SlackChannelList, error)
}

// slackIntegrations implements SlackIntegrations.
type slackIntegrations struct {
	client *Client
}

// SlackEvent represents a event slack integration subscribes to.
type SlackEvent string

const (
	RunApprovalRequiredEvent SlackEvent = "run_approval_required"
	RunSuccessEvent          SlackEvent = "run_success"
	RunErroredEvent          SlackEvent = "run_errored"
)

type SlackStatus string

const (
	IntegrationActive   SlackStatus = "active"
	IntegrationDisabled SlackStatus = "disabled"
)

// SlackIntegration represents a Scalr IACP slack integration.
type SlackIntegration struct {
	ID        string       `jsonapi:"primary,slack-integrations"`
	Name      string       `jsonapi:"attr,name"`
	Status    SlackStatus  `jsonapi:"attr,status"`
	ChannelId string       `jsonapi:"attr,channel-id"`
	Events    []SlackEvent `jsonapi:"attr,events"`

	// Relations
	Environment *Environment  `jsonapi:"relation,environment"`
	Account     *Account      `jsonapi:"relation,account"`
	Workspaces  []*Workspaces `jsonapi:"relation,workspaces"`
}

type SlackIntegrationList struct {
	*Pagination
	Items []*SlackIntegration
}

type SlackIntegrationListOptions struct {
	ListOptions

	Account *string `url:"filter[account]"`
}

type SlackIntegrationCreateOptions struct {
	ID        string       `jsonapi:"primary,slack-integrations"`
	Name      *string      `jsonapi:"attr,name"`
	ChannelId *string      `jsonapi:"attr,channel-id"`
	Events    []SlackEvent `jsonapi:"attr,events"`

	Account     *Account         `jsonapi:"relation,account"`
	Connection  *SlackConnection `jsonapi:"relation,connection"`
	Environment *Environment     `jsonapi:"relation,environment"`
	Workspaces  []*Workspaces    `jsonapi:"relation,workspaces,omitempty"`
}

type SlackIntegrationUpdateOptions struct {
	ID        string       `jsonapi:"primary,slack-integrations"`
	Name      *string      `jsonapi:"attr,name,omitempty"`
	ChannelId *string      `jsonapi:"attr,channel-id,omitempty"`
	Status    SlackStatus  `jsonapi:"attr,status,omitempty"`
	Events    []SlackEvent `jsonapi:"attr,events,omitempty"`

	Environment *Environment  `jsonapi:"relation,environment,omitempty"`
	Workspaces  []*Workspaces `jsonapi:"relation,workspaces,omitempty"`
}

type SlackConnection struct {
	ID                 string `jsonapi:"primary,slack-connections"`
	SlackWorkspaceName string `jsonapi:"attr,slack-workspace-name"`

	// Relations
	Account *Account `jsonapi:"relation,account"`
}

type SlackChannel struct {
	ID        string `json:"slack-connections"`
	Name      string `json:"name"`
	IsPrivate string `json:"is-private"`
}

type SlackChannelList struct {
	*Pagination
	Items []*SlackChannel
}

type SlackChannelListOptions struct {
	ListOptions

	Query *string `url:"query,omitempty"`
}

func (s *slackIntegrations) List(
	ctx context.Context, options SlackIntegrationListOptions,
) (*SlackIntegrationList, error) {
	req, err := s.client.newRequest("GET", "integrations/slack", &options)
	if err != nil {
		return nil, err
	}

	wl := &SlackIntegrationList{}
	err = s.client.do(ctx, req, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

func (s *slackIntegrations) Create(
	ctx context.Context, options SlackIntegrationCreateOptions,
) (*SlackIntegration, error) {
	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "integrations/slack", &options)
	if err != nil {
		return nil, err
	}

	w := &SlackIntegration{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *slackIntegrations) Read(ctx context.Context, si string) (*SlackIntegration, error) {
	if !validStringID(&si) {
		return nil, errors.New("invalid value for Slack integration ID")
	}

	u := fmt.Sprintf("integrations/slack/%s", url.QueryEscape(si))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	w := &SlackIntegration{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *slackIntegrations) Update(
	ctx context.Context, si string, options SlackIntegrationUpdateOptions,
) (*SlackIntegration, error) {
	if !validStringID(&si) {
		return nil, errors.New("invalid value for slack integration ID")
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("integrations/slack/%s", url.QueryEscape(si))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &SlackIntegration{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *slackIntegrations) Delete(ctx context.Context, si string) error {
	if !validStringID(&si) {
		return errors.New("invalid value for slack integration ID")
	}

	u := fmt.Sprintf("integrations/slack/%s", url.QueryEscape(si))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

func (s *slackIntegrations) GetConnection(ctx context.Context, accID string) (*SlackConnection, error) {
	if !validStringID(&accID) {
		return nil, errors.New("invalid value for account ID")
	}

	u := fmt.Sprintf("/integrations/slack/%s/connection", url.QueryEscape(accID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	c := &SlackConnection{}
	err = s.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *slackIntegrations) GetChannels(
	ctx context.Context, accID string, options SlackChannelListOptions,
) (*SlackChannelList, error) {
	if !validStringID(&accID) {
		return nil, errors.New("invalid value for account ID")
	}

	u := fmt.Sprintf("/integrations/slack/%s/connection/channels", url.QueryEscape(accID))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	c := &SlackChannelList{}
	err = s.client.do(ctx, req, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
