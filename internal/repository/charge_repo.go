package repository

import (
	"github.com/diother/go-invoices/internal/models"
	"github.com/jmoiron/sqlx"
)

type ChargeRepository struct {
	db *sqlx.DB
}

func NewChargeRepository(db *sqlx.DB) *ChargeRepository {
	return &ChargeRepository{db: db}
}

func (r *ChargeRepository) Insert(charge models.Charge) error {
	query := `
    INSERT INTO charges (id, created, gross, fee, net, product, client_name, client_email)
    VALUES (:id, :created, :gross, :fee, :net, :product, :client_name, :client_email)
    `
	_, err := r.db.NamedExec(query, charge)
	return err
}
