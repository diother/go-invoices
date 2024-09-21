package services

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/models"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
)

type DonationService struct {
	repo WebhookRepository
}

func NewDonationService(repo WebhookRepository) *DonationService {
	return &DonationService{repo: repo}
}

func (d *DonationService) ProcessDonation(charge *stripe.Charge) (err error) {
	if err = validateCharge(charge); err != nil {
		return fmt.Errorf("Charge validation error: %w", err)
	}

	transaction, err := fetchTransaction(charge.BalanceTransaction.ID)
	if err != nil {
		return fmt.Errorf("Transaction fetch error: %w", err)
	}

	if err = validateChargeTransaction(transaction); err != nil {
		return fmt.Errorf("Transaction validation error: %w", err)
	}

	donation := models.NewDonation(
		transaction.ID,
		uint64(transaction.Created),
		uint32(transaction.Amount),
		uint32(transaction.Fee),
		uint32(transaction.Net),
		charge.BillingDetails.Name,
		charge.BillingDetails.Email,
		sql.NullString{Valid: false},
	)
	if err = d.repo.InsertDonation(donation); err != nil {
		return fmt.Errorf("Database donation insertion failed: %w", err)
	}
	return
}

func fetchTransaction(id string) (*stripe.BalanceTransaction, error) {
	params := &stripe.BalanceTransactionParams{}
	transaction, err := balancetransaction.Get(id, params)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
