package scalr

type Agent struct {
	ID   string `jsonapi:"primary,agents"`
	Name string `jsonapi:"attr,name"`
	OS   string `jsonapi:"attr,os"`
}
