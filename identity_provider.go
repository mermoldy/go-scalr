package scalr

const defaultIdentityProviderLdapID = "idp-sojhv9e8mc2k808"

// IdentityProvider represents a Scalr identity provider.
type IdentityProvider struct {
	ID string `jsonapi:"primary,identity-providers"`
}
