package models

type Payout struct {
	ID        string      `db:"id" json:"id"`
	Created   uint64      `db:"created" json:"created"`
	Gross     uint32      `db:"gross" json:"gross"`
	Fee       uint32      `db:"fee" json:"fee"`
	Net       uint32      `db:"net" json:"net"`
	Donations []*Donation `db:"charges" json:"charges"`
}
