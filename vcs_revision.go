package scalr

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ VcsRevisions = (*vcs_revisions)(nil)

// workspaces implements Workspaces.
type vcs_revisions struct {
	client *Client
}

// VcsRevisions describes all the vcs revisions related methods that the Scalr API supports.
type VcsRevisions interface {
	// Read reads a vcs revisions by its ID.
	Read(ctx context.Context, vcsRevisionId string) (*VcsRevision, error)
}

// VcsRevision represents the VCS metadata
type VcsRevision struct {
	ID             string `jsonapi:"primary,vcs-revisions"`
	Branch         string `jsonapi:"attr,branch"`
	CommitSha      string `jsonapi:"attr,commit-sha"`
	CommitMessage  string `jsonapi:"attr,commit-message"`
	SenderUsername string `jsonapi:"attr,sender-username"`
}

// Read a workspace by its name.
func (s *vcs_revisions) Read(ctx context.Context, vcsRevisionId string) (*VcsRevision, error) {
	options := struct {
		Include string `url:"include"`
	}{
		Include: "created-by",
	}

	u := fmt.Sprintf("vcs-revisions", url.QueryEscape(vcsRevisionId))
	req, err := s.client.newRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	r := &VcsRevision{}
	err = s.client.do(ctx, req, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
