package scalr

type ServiceAccount struct {
	ID      string   `jsonapi:"primary,service-accounts"`
	Name    string   `jsonapi:"attr,name"`
	Email   string   `jsonapi:"attr,email"`
	Status  string   `jsonapi:"attr,status"`
	Account *Account `jsonapi:"relation,account"`
}
