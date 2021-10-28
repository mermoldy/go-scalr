package scalr

// IdentityProvider represents a Scalr identity provider.
type IdentityProvider struct {
	ID string `jsonapi:"primary,identity-providers"`
}
