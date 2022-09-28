package scalr

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ EnvironmentTags = (*environmentTag)(nil)

// EnvironmentTags describes all the environment tags related methods that the
// Scalr API supports.
type EnvironmentTags interface {
	Add(ctx context.Context, envID string, tags []*TagRelation) error
	Replace(ctx context.Context, envID string, tags []*TagRelation) error
	Delete(ctx context.Context, envID string, tags []*TagRelation) error
}

// environmentTag implements EnvironmentTags.
type environmentTag struct {
	client *Client
}

// Add tags to the environment
func (s *environmentTag) Add(ctx context.Context, envID string, trs []*TagRelation) error {
	u := fmt.Sprintf("environments/%s/relationships/tags", url.QueryEscape(envID))
	req, err := s.client.newRequest("POST", u, trs)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Replace environment's tags
func (s *environmentTag) Replace(ctx context.Context, envID string, trs []*TagRelation) error {
	u := fmt.Sprintf("environments/%s/relationships/tags", url.QueryEscape(envID))
	req, err := s.client.newRequest("PATCH", u, trs)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// Delete environment's tags
func (s *environmentTag) Delete(ctx context.Context, envID string, trs []*TagRelation) error {
	u := fmt.Sprintf("environments/%s/relationships/tags", url.QueryEscape(envID))
	req, err := s.client.newRequest("DELETE", u, trs)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
