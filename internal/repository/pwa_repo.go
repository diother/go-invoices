package repository

import (
	"github.com/jmoiron/sqlx"
)

type PWARepository struct {
	db *sqlx.DB
}

func NewPWARepository(db *sqlx.DB) *PWARepository {
	return &PWARepository{db: db}
}
