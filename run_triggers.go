package scalr

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ RunTriggers = (*runTriggers)(nil)

type RunTriggers interface {
	// Create is used to create a new run trigger.
	Create(ctx context.Context, options RunTriggerCreateOptions) (*RunTrigger, error)

	// Read RunTrigger by it's ID
	Read(ctx context.Context, runTriggerID string) (*RunTrigger, error)

	// Delete RunTrigger by it's ID
	Delete(ctx context.Context, runTriggerID string) error
}

// runTriggers implements RunTriggers.
type runTriggers struct {
	client *Client
}

type RunTrigger struct {
	ID        string    `jsonapi:"primary,run-triggers"`
	CreatedAt time.Time `jsonapi:"attr,created-at,iso8601"`

	// Relations
	Upstream   *Upstream   `jsonapi:"relation,upstream"`
	Downstream *Downstream `jsonapi:"relation,downstream"`
}

type RunTriggerCreateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,run-triggers"`

	Downstream *Downstream `jsonapi:"relation,downstream"`
	Upstream   *Upstream   `jsonapi:"relation,upstream"`
}

type Downstream struct {
	ID string `jsonapi:"primary,workspaces"`
}

type Upstream struct {
	ID string `jsonapi:"primary,workspaces"`
}

// Create is used to create a new runTrigger.
func (s *runTriggers) Create(ctx context.Context, options RunTriggerCreateOptions) (*RunTrigger, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	// Make sure we don't send a user provided ID.
	options.ID = ""

	req, err := s.client.newRequest("POST", "run-triggers", &options)
	if err != nil {
		return nil, err
	}

	runTrigger := &RunTrigger{}
	err = s.client.do(ctx, req, runTrigger)
	if err != nil {
		return nil, err
	}

	return runTrigger, nil
}

func (o RunTriggerCreateOptions) valid() error {
	if o.Downstream == nil {
		return errors.New("downstream ID is required")
	}
	if o.Upstream == nil {
		return errors.New("upstream ID is required")
	}
	if !validString(&o.Downstream.ID) {
		return errors.New("downstream ID is required")
	}
	if !validStringID(&o.Downstream.ID) {
		return errors.New("invalid value for Downstream ID")
	}
	if !validString(&o.Upstream.ID) {
		return errors.New("upstream ID is required")
	}
	if !validStringID(&o.Upstream.ID) {
		return errors.New("invalid value for Upstream ID")
	}
	return nil
}

func (s *runTriggers) Read(ctx context.Context, runTriggerID string) (*RunTrigger, error) {
	if !validStringID(&runTriggerID) {
		return nil, errors.New("invalid value for RunTrigger ID")
	}
	u := fmt.Sprintf("run-triggers/%s", url.QueryEscape(runTriggerID))
	fmt.Println(u)
	req, err := s.client.newRequest("GET", u, nil)

	if err != nil {
		return nil, err
	}

	runTrigger := &RunTrigger{}
	err = s.client.do(ctx, req, runTrigger)
	if err != nil {
		return nil, err
	}

	return runTrigger, nil
}

func (s *runTriggers) Delete(ctx context.Context, runTriggerID string) error {
	if !validStringID(&runTriggerID) {
		return errors.New("invalid value for RunTrigger ID")
	}
	u := fmt.Sprintf("run-triggers/%s", url.QueryEscape(runTriggerID))
	fmt.Println(u)
	req, err := s.client.newRequest("DELETE", u, nil)

	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
