package scalr

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ AgentPoolTokens = (*agentPoolTokens)(nil)

// AgentPoolTokens describes all the access token related methods that the
// Scalr IACP API supports.
type AgentPoolTokens interface {
	List(ctx context.Context, agentPoolID string, options AccessTokenListOptions) (*AccessTokenList, error)
	Create(ctx context.Context, agentPoolID string, options AccessTokenCreateOptions) (*AccessToken, error)
}

// agentPoolTokens implements AgentPoolTokens.
type agentPoolTokens struct {
	client *Client
}

// List all the agent pool's tokens.
func (s *agentPoolTokens) List(ctx context.Context, agentPoolID string, options AccessTokenListOptions) (*AccessTokenList, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("agent-pools/%s/access-tokens", url.QueryEscape(agentPoolID)), &options)
	if err != nil {
		return nil, err
	}

	tl := &AccessTokenList{}
	err = s.client.do(ctx, req, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// Create is used to create a new AccessToken for AgentPool.
func (s *agentPoolTokens) Create(ctx context.Context, agentPoolID string, options AccessTokenCreateOptions) (*AccessToken, error) {

	// Make sure we don't send a user provided ID.
	options.ID = ""

	if !validStringID(&agentPoolID) {
		return nil, fmt.Errorf("invalid value for agent pool ID: '%s'", agentPoolID)
	}

	req, err := s.client.newRequest("POST", fmt.Sprintf("agent-pools/%s/access-tokens", url.QueryEscape(agentPoolID)), &options)
	if err != nil {
		return nil, err
	}

	agentPoolToken := &AccessToken{}
	err = s.client.do(ctx, req, agentPoolToken)
	if err != nil {
		return nil, err
	}

	return agentPoolToken, nil
}
