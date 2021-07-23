package scalr

type Team struct {
	ID   string `jsonapi:"primary,teams"`
	Name string `jsonapi:"attr,name,omitempty"`
}
