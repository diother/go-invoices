package repository

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/custom_errors"
	"github.com/diother/go-invoices/internal/models"
)

func (r *AuthRepository) GetUser(username string) (*models.User, error) {
	var donation models.User
	query := "SELECT * FROM users WHERE username = ?"

	if err := r.db.Get(&donation, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NewCredentialsError("username or password is invalid")
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return &donation, nil
}
