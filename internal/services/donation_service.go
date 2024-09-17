package services

import (
	"github.com/diother/go-invoices/internal/models"
)

type DonationRepository interface {
	Insert(*models.Donation) error
}

type DonationServiceImpl struct {
	repo DonationRepository
}

func NewDonationServiceImpl(repo DonationRepository) *DonationServiceImpl {
	return &DonationServiceImpl{repo: repo}
}
