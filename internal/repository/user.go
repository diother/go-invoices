package repository

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/custom_errors"
	"github.com/diother/go-invoices/internal/models"
)

func (r *AuthRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE username = ?"

	if err := r.db.Get(&user, query, username); err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NewCredentialsError("Nume de utilizator sau parolă invalide")
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = ?"

	if err := r.db.Get(&user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with the id: %v", id)
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return &user, nil
}
