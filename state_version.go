package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ StateVersions = (*state_versions)(nil)

// StateVersions describes all the state versions related methods that the Scalr API supports.
type StateVersions interface {
	ReadByID(ctx context.Context, stateVersionID string) (*StateVersion, error)
	ReadCurrentFromWorkspace(ctx context.Context, workspaceID string) (*StateVersion, error)
}

type StateVersion struct {
	ID        string      `jsonapi:"primary,state-versions"`
	CreatedAt time.Time   `jsonapi:"attr,created-at,iso8601"`
	Force     bool        `jsonapi:"attr,force"`
	Resources []*Resource `jsonapi:"attr,resources"`

	// Relations
	Run       *Run       `jsonapi:"relation,run"`
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

type Resource struct {
	Type string `jsonapi:"attr,type"`
}

// state_versions implements StateVersion.
type state_versions struct {
	client *Client
}

// Read a state version by its ID.
func (s *state_versions) ReadByID(ctx context.Context, stateVersionID string) (*StateVersion, error) {
	if !validStringID(&stateVersionID) {
		return nil, errors.New("invalid value for state version")
	}

	u := fmt.Sprintf("state-versions/%s", url.QueryEscape(stateVersionID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = s.client.do(ctx, req, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Read current state version of workspace.
func (s *state_versions) ReadCurrentFromWorkspace(ctx context.Context, workspaceID string) (*StateVersion, error) {
	if !validStringID(&workspaceID) {
		return nil, errors.New("invalid value for workspace")
	}

	u := fmt.Sprintf("workspaces/%s/current-state-version", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = s.client.do(ctx, req, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}
