package models

type Donation struct {
	ID          string `db:"id" json:"id"`
	Created     uint64 `db:"created" json:"created"`
	Gross       uint32 `db:"gross" json:"gross"`
	Fee         uint32 `db:"fee" json:"fee"`
	Net         uint32 `db:"net" json:"net"`
	Product     string `db:"product" json:"product"`
	ClientName  string `db:"client_name" json:"client_name"`
	ClientEmail string `db:"client_email" json:"client_email"`
}
