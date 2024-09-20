package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type WebhookRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewWebhookRepository(db *sqlx.DB) *WebhookRepository {
	return &WebhookRepository{db: db}
}

func (r *WebhookRepository) BeginTransaction() error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	r.tx = tx
	return nil
}

func (r *WebhookRepository) Commit() error {
	if r.tx == nil {
		return fmt.Errorf("No transaction to commit")
	}
	err := r.tx.Commit()
	r.tx = nil
	return err
}

func (r *WebhookRepository) Rollback() error {
	if r.tx == nil {
		return fmt.Errorf("No transaction to rollback")
	}
	err := r.tx.Rollback()
	r.tx = nil
	return err
}

func (r *WebhookRepository) execNamed(query string, arg interface{}) (sql.Result, error) {
	if r.tx != nil {
		return r.tx.NamedExec(query, arg)
	}
	return r.db.NamedExec(query, arg)
}
