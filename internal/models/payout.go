package models

type Payout struct {
	ID        string      `db:"id"`
	Created   uint64      `db:"created"`
	Gross     uint32      `db:"gross"`
	Fee       uint32      `db:"fee"`
	Net       uint32      `db:"net"`
	Donations []*Donation `db:"charges"`
}
