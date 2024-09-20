package repository

import (
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
