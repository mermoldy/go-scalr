package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ WorkspaceTags = (*workspaceTag)(nil)

// WorkspaceTags describes all the workspace tags related methods that the
// Scalr API supports.
type WorkspaceTags interface {
	Create(ctx context.Context, options WorkspaceTagsCreateOptions) error
	Update(ctx context.Context, options WorkspaceTagsUpdateOptions) error
}

// workspaceTag implements WorkspaceTags.
type workspaceTag struct {
	client *Client
}

// WorkspaceTag represents a single workspace tag relation.
type WorkspaceTag struct {
	ID string `jsonapi:"primary,tags"`
}

// WorkspaceTagsCreateOptions represents options for adding tags to a workspace.
type WorkspaceTagsCreateOptions struct {
	WorkspaceID   string
	WorkspaceTags []*WorkspaceTag
}

// WorkspaceTagsUpdateOptions represents options for updating tags in a workspace.
type WorkspaceTagsUpdateOptions struct {
	WorkspaceID   string
	WorkspaceTags []*WorkspaceTag
}

func (o WorkspaceTagsCreateOptions) valid() error {
	if !validStringID(&o.WorkspaceID) {
		return errors.New("invalid value for workspace ID")
	}
	if o.WorkspaceTags == nil || len(o.WorkspaceTags) < 1 {
		return errors.New("list of tags is required")
	}
	return nil
}

func (o WorkspaceTagsUpdateOptions) valid() error {
	if !validStringID(&o.WorkspaceID) {
		return errors.New("invalid value for workspace ID")
	}
	return nil
}

// Create is used for adding tags to the workspace.
func (s *workspaceTag) Create(ctx context.Context, options WorkspaceTagsCreateOptions) error {
	if err := options.valid(); err != nil {
		return err
	}
	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.QueryEscape(options.WorkspaceID))
	req, err := s.client.newRequest("POST", u, options.WorkspaceTags)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Update is used for tags replacement in the workspace.
func (s *workspaceTag) Update(ctx context.Context, options WorkspaceTagsUpdateOptions) error {
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.QueryEscape(options.WorkspaceID))
	req, err := s.client.newRequest("PATCH", u, options.WorkspaceTags)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
