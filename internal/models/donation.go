package models

type Donation struct {
	ID          string `db:"id"`
	Created     uint64 `db:"created"`
	Gross       uint32 `db:"gross"`
	Fee         uint32 `db:"fee"`
	Net         uint32 `db:"net"`
	ClientName  string `db:"client_name"`
	ClientEmail string `db:"client_email"`
	PayoutID    string `db:"payout_id"`
}
