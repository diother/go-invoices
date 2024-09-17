package repository

import (
	"github.com/diother/go-invoices/internal/models"
	"github.com/jmoiron/sqlx"
)

type PayoutRepositoryMySQL struct {
	db *sqlx.DB
}

func NewPayoutRepositoryMySQL(db *sqlx.DB) *PayoutRepositoryMySQL {
	return &PayoutRepositoryMySQL{db: db}
}

func (r *PayoutRepositoryMySQL) Insert(payout *models.Payout) error {
	query := `
    INSERT INTO payouts (id, created, gross, fee, net)
    VALUES (:id, :created, :gross, :fee, :net)
    `
	_, err := r.db.NamedExec(query, payout)
	return err
}
