package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ Tags = (*tags)(nil)

// Tags describes all the tags related methods that the Scalr API supports.
type Tags interface {
	// List all the tags.
	List(ctx context.Context, options TagListOptions) (*TagList, error)
	// Create is used to create a new tag.
	Create(ctx context.Context, options TagCreateOptions) (*Tag, error)
	// Read reads a tag by its ID.
	Read(ctx context.Context, tagID string) (*Tag, error)
	// Update existing tag by its ID.
	Update(ctx context.Context, tagID string, options TagUpdateOptions) (*Tag, error)
	// Delete deletes a tag by its ID.
	Delete(ctx context.Context, tagID string) error
}

// tags implements Tags.
type tags struct {
	client *Client
}

// TagList represents a list of tags.
type TagList struct {
	*Pagination
	Items []*Tag
}

type Tag struct {
	ID   string `jsonapi:"primary,tags"`
	Name string `jsonapi:"attr,name"`

	// Relations
	Account *Account `jsonapi:"relation,account"`
}

type TagRelation struct {
	ID string `jsonapi:"primary,tags"`
}

// TagListOptions represents the options for listing tags.
type TagListOptions struct {
	ListOptions

	ID      *string `url:"filter[tag],omitempty"`
	Account *string `url:"filter[account],omitempty"`
	Name    *string `url:"filter[name],omitempty"`
	Query   *string `url:"query,omitempty"`
}

// TagCreateOptions represents the options for creating a new tag.
type TagCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,tags"`
	// The name of the tag, it must be unique within the account.
	Name *string `jsonapi:"attr,name"`
	// Specifies the Account for the tag.
	Account *Account `jsonapi:"relation,account"`
}

// TagUpdateOptions represents the options for updating a tag.
type TagUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,tags"`
	// The name of the tag, it must be unique within the account.
	Name *string `jsonapi:"attr,name"`
}

// List all the tags.
func (s *tags) List(ctx context.Context, options TagListOptions) (*TagList, error) {
	req, err := s.client.newRequest("GET", "tags", &options)
	if err != nil {
		return nil, err
	}

	tl := &TagList{}
	err = s.client.do(ctx, req, tl)
	if err != nil {
		return nil, err
	}

	return tl, nil
}

// Read reads a tag by its ID.
func (s *tags) Read(ctx context.Context, tagID string) (*Tag, error) {
	if !validStringID(&tagID) {
		return nil, errors.New("invalid value for tag ID")
	}

	u := fmt.Sprintf("tags/%s", url.QueryEscape(tagID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	t := &Tag{}
	err = s.client.do(ctx, req, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (o TagCreateOptions) valid() error {
	if o.Account == nil {
		return errors.New("account is required")
	}
	if !validStringID(&o.Account.ID) {
		return errors.New("invalid value for account ID")
	}
	if o.Name == nil {
		return errors.New("name is required")
	}
	return nil
}

// Create is used to create a new tag.
func (s *tags) Create(ctx context.Context, options TagCreateOptions) (*Tag, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}
	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "tags", &options)
	if err != nil {
		return nil, err
	}

	t := &Tag{}
	err = s.client.do(ctx, req, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Update is used to update a tag.
func (s *tags) Update(ctx context.Context, tagID string, options TagUpdateOptions) (*Tag, error) {
	if !validStringID(&tagID) {
		return nil, errors.New("invalid value for tag ID")
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	u := fmt.Sprintf("tags/%s", url.QueryEscape(tagID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	t := &Tag{}
	err = s.client.do(ctx, req, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Delete tag by its ID.
func (s *tags) Delete(ctx context.Context, tagID string) error {
	if !validStringID(&tagID) {
		return errors.New("invalid value for tag ID")
	}

	u := fmt.Sprintf("tags/%s", url.QueryEscape(tagID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
