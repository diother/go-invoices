package repository

import (
	"github.com/diother/go-invoices/internal/models"
	"github.com/jmoiron/sqlx"
)

type DonationRepositoryMySQL struct {
	db *sqlx.DB
}

func NewDonationRepositoryMySQL(db *sqlx.DB) *DonationRepositoryMySQL {
	return &DonationRepositoryMySQL{db: db}
}

func (r *DonationRepositoryMySQL) Insert(donation *models.Donation) error {
	query := `
    INSERT INTO donations (id, created, gross, fee, net, product, client_name, client_email)
    VALUES (:id, :created, :gross, :fee, :net, :product, :client_name, :client_email)
    `
	_, err := r.db.NamedExec(query, donation)
	return err
}
