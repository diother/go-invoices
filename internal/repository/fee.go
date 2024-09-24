package repository

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/models"
)

func (r *WebhookRepository) InsertFee(fee *models.Fee) error {
	query := `
    INSERT INTO fees (id, description, created, fee, payout_id)
	VALUES (:id, :description, :created, :fee, :payout_id)
    `
	_, err := r.execNamed(query, fee)
	return err
}

func (r *PWARepository) GetRelatedFees(payoutID string) (fees []*models.Fee, err error) {
	query := "SELECT id, description, created, fee FROM fees WHERE payout_id = ?"

	if err := r.db.Select(&fees, query, payoutID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no fees with payout_id: %s", payoutID)
		}
		return nil, fmt.Errorf("failed to retrieve fees: %w", err)
	}
	return
}
