package repository

import (
	"github.com/diother/go-invoices/internal/models"
	"github.com/jmoiron/sqlx"
)

type PayoutRepository struct {
	db *sqlx.DB
}

func NewPayoutRepository(db *sqlx.DB) *PayoutRepository {
	return &PayoutRepository{db: db}
}

func (r *PayoutRepository) Insert(payout models.Payout) error {
	query := `
    INSERT INTO payouts (id, created, gross, fee, net)
    VALUES (:id, :created, :gross, :fee, :net)
    `
	_, err := r.db.NamedExec(query, payout)
	return err
}
