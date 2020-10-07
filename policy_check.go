package scalr

import (
	"time"
)

// PolicyScope represents a policy scope.
type PolicyScope string

// List all available policy scopes.
const (
	PolicyScopeEnvironment PolicyScope = "environment"
	PolicyScopeWorkspace   PolicyScope = "workspace"
)

// PolicyStatus represents a policy check state.
type PolicyStatus string

//List all available policy check statuses.
const (
	PolicyCanceled    PolicyStatus = "canceled"
	PolicyErrored     PolicyStatus = "errored"
	PolicyHardFailed  PolicyStatus = "hard_failed"
	PolicyOverridden  PolicyStatus = "overridden"
	PolicyPasses      PolicyStatus = "passed"
	PolicyPending     PolicyStatus = "pending"
	PolicyQueued      PolicyStatus = "queued"
	PolicySoftFailed  PolicyStatus = "soft_failed"
	PolicyUnreachable PolicyStatus = "unreachable"
)

// PolicyCheck represents a Scalr policy check..
type PolicyCheck struct {
	ID               string                  `jsonapi:"primary,policy-checks"`
	Actions          *PolicyActions          `jsonapi:"attr,actions"`
	Permissions      *PolicyPermissions      `jsonapi:"attr,permissions"`
	Result           *PolicyResult           `jsonapi:"attr,result"`
	Scope            PolicyScope             `jsonapi:"attr,scope"`
	Status           PolicyStatus            `jsonapi:"attr,status"`
	StatusTimestamps *PolicyStatusTimestamps `jsonapi:"attr,status-timestamps"`
}

// PolicyActions represents the policy check actions.
type PolicyActions struct {
	IsOverridable bool `json:"is-overridable"`
}

// PolicyPermissions represents the policy check permissions.
type PolicyPermissions struct {
	CanOverride bool `json:"can-override"`
}

// PolicyResult represents the complete policy check result,
type PolicyResult struct {
	AdvisoryFailed int  `json:"advisory-failed"`
	Duration       int  `json:"duration"`
	HardFailed     int  `json:"hard-failed"`
	Passed         int  `json:"passed"`
	Result         bool `json:"result"`
	SoftFailed     int  `json:"soft-failed"`
	TotalFailed    int  `json:"total-failed"`
}

// PolicyStatusTimestamps holds the timestamps for individual policy check
// statuses.
type PolicyStatusTimestamps struct {
	ErroredAt    time.Time `json:"errored-at"`
	HardFailedAt time.Time `json:"hard-failed-at"`
	PassedAt     time.Time `json:"passed-at"`
	QueuedAt     time.Time `json:"queued-at"`
	SoftFailedAt time.Time `json:"soft-failed-at"`
}
