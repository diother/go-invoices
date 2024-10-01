package repository

import "github.com/diother/go-invoices/internal/models"

func (r *AuthRepository) InsertSession(payout *models.Session) error {
	query := `
    INSERT INTO sessions (session_token, user_id, expires_at)
	VALUES (:session_token, :user_id, :expires_at)
    `
	_, err := r.db.NamedExec(query, payout)
	return err
}
