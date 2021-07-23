package scalr

type ServiceAccount struct {
	ID      string   `jsonapi:"primary,service-accounts"`
	Name    string   `jsonapi:"attr,name,omitempty"`
	Email   string   `jsonapi:"attr,email,omitempty"`
	Status  string   `jsonapi:"attr,status,omitempty"`
	Account *Account `jsonapi:"relation,account,omitempty"`
}
