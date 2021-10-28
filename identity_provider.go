package scalr

const (
	defaultIdentityProviderScalrID = "idp-sohkb0o1phrdmr8"
	defaultIdentityProviderLdapID  = "idp-sojhv9e8mc2k808"
)

// IdentityProvider represents a Scalr identity provider.
type IdentityProvider struct {
	ID string `jsonapi:"primary,identity-providers"`
}
