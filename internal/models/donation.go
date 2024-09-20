package models

import "database/sql"

type Donation struct {
	ID          string         `db:"id"`
	Created     uint64         `db:"created"`
	Gross       uint32         `db:"gross"`
	Fee         uint32         `db:"fee"`
	Net         uint32         `db:"net"`
	ClientName  string         `db:"client_name"`
	ClientEmail string         `db:"client_email"`
	PayoutID    sql.NullString `db:"payout_id"`
}

func NewDonation(id string, created uint64, gross, fee, net uint32, clientName, clientEmail string, payoutID sql.NullString) *Donation {
	return &Donation{
		ID:          id,
		Created:     created,
		Gross:       gross,
		Fee:         fee,
		Net:         net,
		ClientName:  clientName,
		ClientEmail: clientEmail,
		PayoutID:    payoutID,
	}
}
