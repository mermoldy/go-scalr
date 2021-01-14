package scalr

// User represents a Scalr user.
type User struct {
	ID       string `jsonapi:"primary,users"`
	Email    string `jsonapi:"attr,email"`
	Username string `jsonapi:"attr,username"`
	FullName string `jsonapi:"attr,full-name"`
}
