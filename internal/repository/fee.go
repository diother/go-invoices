package repository

import "github.com/diother/go-invoices/internal/models"

func (r *WebhookRepository) InsertFee(fee *models.Fee) error {
	query := `
    INSERT INTO fees (id, created, fee, payout_id)
	VALUES (:id, :created, :fee, :payout_id)
    `
	_, err := r.execNamed(query, fee)
	return err
}
