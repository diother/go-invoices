package repository

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/models"
)

func (r *WebhookRepository) InsertDonation(donation *models.Donation) error {
	query := `
    INSERT INTO donations (id, created, gross, fee, net, client_name, client_email, payout_id)
	VALUES (:id, :created, :gross, :fee, :net, :client_name, :client_email, :payout_id)
    `
	_, err := r.execNamed(query, donation)
	return err
}

func (r *WebhookRepository) UpdateRelatedPayout(donation *models.Donation) (bool, error) {
	query := `
	UPDATE donations
	SET payout_id = :payout_id
	WHERE id = :id
	`
	result, err := r.execNamed(query, donation)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected != 0, nil
}

func (r *PWARepository) GetAllDonations() ([]*models.Donation, error) {
	var donations []*models.Donation
	query := "SELECT * FROM donations"

	if err := r.db.Select(&donations, query); err != nil {
		return nil, err
	}
	return donations, nil
}

func (r *PWARepository) GetDonation(id string) (*models.Donation, error) {
	var donation models.Donation
	query := "SELECT * FROM donations WHERE id = ?"

	if err := r.db.Get(&donation, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("donation with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve donation: %w", err)
	}
	return &donation, nil
}
