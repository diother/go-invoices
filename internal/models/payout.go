package models

type Payout struct {
	ID      string `db:"id"`
	Created uint64 `db:"created"`
	Gross   uint32 `db:"gross"`
	Fee     uint32 `db:"fee"`
	Net     uint32 `db:"net"`
}

func NewPayout(id string, created uint64, gross, fee, net uint32) *Payout {
	return &Payout{
		ID:      id,
		Created: created,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}
