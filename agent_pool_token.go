package scalr

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ AgentPoolTokens = (*agentPoolTokens)(nil)

// AgentPoolTokens describes all the access token related methods that the
// Scalr IACP API supports.
type AgentPoolTokens interface {
	List(ctx context.Context, agentPoolID string) (*AgentPoolTokenList, error)
	Create(ctx context.Context, agentPoolID string, options AgentPoolTokenCreateOptions) (*AgentPoolToken, error)
}

// agentPoolTokens implements AgentPoolTokens.
type agentPoolTokens struct {
	client *Client
}

// AgentPoolTokenList represents a list of agent pools.
type AgentPoolTokenList struct {
	*Pagination
	Items []*AgentPoolToken
}

// AgentPoolToken represents a Scalr agent pool.
type AgentPoolToken struct {
	ID          string    `jsonapi:"primary,access-tokens"`
	CreatedAt   time.Time `jsonapi:"attr,created-at,iso8601"`
	Description string    `jsonapi:"attr,description"`
	Token       string    `jsonapi:"attr,token"`
}

// AgentPoolTokenCreateOptions represents the options for creating a new AgentPoolToken.
type AgentPoolTokenCreateOptions struct {
	ID          string  `jsonapi:"primary,access-tokens"`
	Description *string `jsonapi:"attr,description,omitempty"`
}

// List all the agent pools.
func (s *agentPoolTokens) List(ctx context.Context, agentPoolID string) (*AgentPoolTokenList, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("agent-pools/%s/access-tokens", url.QueryEscape(agentPoolID)), nil)
	if err != nil {
		return nil, err
	}

	tl := &AgentPoolTokenList{}
	err = s.client.do(ctx, req, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// Create is used to create a new AgentPoolToken.
func (s *agentPoolTokens) Create(ctx context.Context, agentPoolID string, options AgentPoolTokenCreateOptions) (*AgentPoolToken, error) {

	// Make sure we don't send a user provided ID.
	options.ID = ""

	if !validStringID(&agentPoolID) {
		return nil, fmt.Errorf("invalid value for agent pool ID: '%s'", agentPoolID)
	}

	req, err := s.client.newRequest("POST", fmt.Sprintf("agent-pools/%s/access-tokens", url.QueryEscape(agentPoolID)), &options)
	if err != nil {
		return nil, err
	}

	agentPoolToken := &AgentPoolToken{}
	err = s.client.do(ctx, req, agentPoolToken)
	if err != nil {
		return nil, err
	}

	return agentPoolToken, nil
}
