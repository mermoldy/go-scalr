package scalr

// VcsProvider represents a Scalr VCS provider.
type VcsProvider struct {
	ID      string `jsonapi:"primary,vcs-providers"`
	VcsType string `jsonapi:"attr,vcs-type"`
	Url     string `jsonapi:"attr,url"`
}
