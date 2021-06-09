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
	Create(ctx context.Context, options StateVersionCreateOptions) (*StateVersion, error)
	ReadByID(ctx context.Context, stateVersionID string) (*StateVersion, error)
	ReadCurrentFromWorkspace(ctx context.Context, workspaceID string) (*StateVersion, error)
}

type StateVersion struct {
	ID        string                 `jsonapi:"primary,state-versions"`
	Force     bool                   `jsonapi:"attr,force"`
	Lineage   string                 `jsonapi:"attr,lineage"`
	MD5       *string                `jsonapi:"attr,md5"`
	CreatedAt time.Time              `jsonapi:"attr,created-at,iso8601"`
	Serial    uint64                 `jsonapi:"attr,serial"`
	Size      uint64                 `jsonapi:"attr,size"`
	State     *string                `jsonapi:"attr,state"`
	Resources []*Resource            `jsonapi:"attr,resources"`
	Outputs   []*Output              `jsonapi:"attr,outputs"`
	Modules   map[string]interface{} `jsonapi:"attr,modules"`
	Providers map[string]interface{} `jsonapi:"attr,providers"`

	Workspace            *Workspace    `jsonapi:"relation,workspace"`
	Run                  *Run          `jsonapi:"relation,run"`
	NextStateVErsion     *StateVersion `jsonapi:"relation,next-state-version"`
	PreviousStateVersion *StateVersion `jsonapi:"relation,previous-state-version"`
}

type Output struct {
	Name      string  `jsonapi:"attr,name"`
	Sensitive bool    `jsonapi:"attr,sensitive"`
	Value     *string `jsonapi:"attr,value"`
}

type Resource struct {
	Type    string  `jsonapi:"attr,type"`
	Address string  `jsonapi:"attr,address"`
	Module  *string `jsonapi:"attr,module"`
}

// state_versions implements StateVersion.
type state_versions struct {
	client *Client
}

type StateVersionCreateOptions struct {
	ID        string                 `jsonapi:"primary,state-versions"`
	Force     bool                   `jsonapi:"attr,force"`
	Lineage   string                 `jsonapi:"attr,lineage"`
	MD5       string                 `jsonapi:"attr,md5"`
	Serial    uint64                 `jsonapi:"attr,serial"`
	Size      uint64                 `jsonapi:"attr,size"`
	State     string                 `jsonapi:"attr,state"`
	Resources []*Resource            `jsonapi:"attr,resources"`
	Outputs   []*Output              `jsonapi:"attr,outputs"`
	Modules   map[string]interface{} `jsonapi:"attr,modules"`
	Providers map[string]interface{} `jsonapi:"attr,providers"`

	Workspace            *Workspace    `jsonapi:"relation,workspace"`
	Run                  *Run          `jsonapi:"relation,run"`
	NextStateVErsion     *StateVersion `jsonapi:"relation,next-state-version"`
	PreviousStateVersion *StateVersion `jsonapi:"relation,previous-state-version"`
}

// Read current state version of workspace.
func (s *state_versions) Create(ctx context.Context, options StateVersionCreateOptions) (*StateVersion, error) {
	options.ID = ""

	req, err := s.client.newRequest("POST", "state-versions", &options)
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
