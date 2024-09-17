package services

import (
	"github.com/diother/go-invoices/internal/models"
)

type PayoutRepository interface {
	Insert(*models.Payout) error
}

type PayoutServiceImpl struct {
	repo PayoutRepository
}

func NewPayoutServiceImpl(repo PayoutRepository) *PayoutServiceImpl {
	return &PayoutServiceImpl{repo: repo}
}
