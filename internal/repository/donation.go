package repository

import (
	"github.com/diother/go-invoices/internal/models"
)

func (r *WebhookRepository) InsertDonation(donation *models.Donation) error {
	query := `
    INSERT INTO donations (id, created, gross, fee, net, client_name, client_email, payout_id)
	VALUES (:id, :created, :gross, :fee, :net, :client_name, :client_email, :payout_id)
    `
	_, err := r.db.NamedExec(query, donation)
	return err
}

func (r *WebhookRepository) UpdateRelatedPayout(donation *models.Donation) (bool, error) {
	query := `
	UPDATE donations
	SET payout_id = :payout_id
	WHERE id = :id
	`
	result, err := r.db.NamedExec(query, donation)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected != 0, nil
}
