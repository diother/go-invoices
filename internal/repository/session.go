package repository

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/models"
)

func (r *AuthRepository) InsertSession(payout *models.Session) error {
	query := `
    INSERT INTO sessions (session_token, user_id, expires_at)
	VALUES (:session_token, :user_id, :expires_at)
    `
	_, err := r.db.NamedExec(query, payout)
	return err
}

func (r *AuthRepository) GetSession(sessionToken int64) (*models.Session, error) {
	var session models.Session
	query := "SELECT * FROM sessions WHERE session_token = ?"

	if err := r.db.Get(&session, query, sessionToken); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no session with the id: %v", sessionToken)
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return &session, nil
}
