package repository

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/models"
)

func (r *WebhookRepository) InsertPayout(payout *models.Payout) error {
	query := `
    INSERT INTO payouts (id, created, gross, fee, net)
    VALUES (:id, :created, :gross, :fee, :net)
    `
	_, err := r.execNamed(query, payout)
	return err
}

func (r *WebhookRepository) GetAllPayouts() ([]*models.Payout, error) {
	var payouts []*models.Payout
	query := "SELECT * FROM payouts"

	if err := r.db.Select(&payouts, query); err != nil {
		return nil, err
	}
	return payouts, nil
}

func (r *WebhookRepository) GetPayout(id string) (*models.Payout, error) {
	var payout models.Payout
	query := "SELECT * FROM payouts WHERE id = ?"

	if err := r.db.Get(&payout, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payout with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve payout: %w", err)
	}
	return &payout, nil
}
