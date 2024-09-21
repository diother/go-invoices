package models

import "database/sql"

type Fee struct {
	ID       string         `db:"id"`
	Created  uint64         `db:"created"`
	Fee      uint32         `db:"fee"`
	PayoutID sql.NullString `db:"payout_id"`
}

func NewFee(id string, created uint64, fee uint32, payoutID sql.NullString) *Fee {
	return &Fee{
		ID:       id,
		Created:  created,
		Fee:      fee,
		PayoutID: payoutID,
	}
}
