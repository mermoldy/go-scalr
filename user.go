package scalr

// User represents a Scalr user.
type User struct {
	ID       string `jsonapi:"primary,users"`
	Email    string `jsonapi:"attr,email,omitempty"`
	Username string `jsonapi:"attr,username,omitempty"`
	FullName string `jsonapi:"attr,full-name,omitempty"`
}
