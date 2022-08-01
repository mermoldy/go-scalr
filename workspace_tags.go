package scalr

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ WorkspaceTags = (*workspaceTag)(nil)

// WorkspaceTags describes all the workspace tags related methods that the
// Scalr API supports.
type WorkspaceTags interface {
	Add(ctx context.Context, wsID string, tags []*TagRelation) error
	Replace(ctx context.Context, wsID string, tags []*TagRelation) error
	Delete(ctx context.Context, wsID string, tags []*TagRelation) error
}

// workspaceTag implements WorkspaceTags.
type workspaceTag struct {
	client *Client
}

// Add tags to the workspace
func (s *workspaceTag) Add(ctx context.Context, wsID string, trs []*TagRelation) error {
	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.QueryEscape(wsID))
	req, err := s.client.newRequest("POST", u, trs)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Replace workspace's tags
func (s *workspaceTag) Replace(ctx context.Context, wsID string, trs []*TagRelation) error {
	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.QueryEscape(wsID))
	req, err := s.client.newRequest("PATCH", u, trs)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Delete workspace's tags
func (s *workspaceTag) Delete(ctx context.Context, wsID string, trs []*TagRelation) error {
	u := fmt.Sprintf("workspaces/%s/relationships/tags", url.QueryEscape(wsID))
	req, err := s.client.newRequest("DELETE", u, trs)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
