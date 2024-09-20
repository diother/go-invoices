package services

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/errors"
	"github.com/diother/go-invoices/internal/models"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
	"github.com/stripe/stripe-go/v79/charge"
)

type WebhookRepository interface {
	InsertDonation(donation *models.Donation) error
	InsertPayout(payout *models.Payout) error
	UpdateRelatedPayout(donation *models.Donation) (bool, error)
	BeginTransaction() error
	Rollback() error
	Commit() error
}

type PayoutService struct {
	repo WebhookRepository
}

func NewPayoutService(repo WebhookRepository) *PayoutService {
	return &PayoutService{repo: repo}
}

func (p *PayoutService) ProcessPayout(payout *stripe.Payout) error {
	if err := validatePayout(payout); err != nil {
		return fmt.Errorf("Payout validation error: %w", err)
	}

	transactions, err := fetchTransactions(payout.ID)
	if err != nil {
		return fmt.Errorf("Transactions fetch failed: %w", err)
	}

	if err = validatePayoutTransaction(transactions[0]); err != nil {
		return fmt.Errorf("Payout transaction validation failed: %w", err)
	}

	for i := 1; i < len(transactions); i++ {
		if err = validateChargeTransaction(transactions[i]); err != nil {
			return fmt.Errorf("Charge transaction validation failed for %s: %w", transactions[i].ID, err)
		}
	}

	payoutGross, payoutFee, payoutNet, err := validateMatchingSums(transactions)
	if err != nil {
		return fmt.Errorf("Matching sum validation failed: %w", err)
	}

	payoutModel := models.NewPayout(
		transactions[0].ID,
		uint64(transactions[0].Created),
		uint32(payoutGross),
		uint32(payoutFee),
		uint32(payoutNet),
	)

	if err := p.repo.BeginTransaction(); err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			p.repo.Rollback()
		} else {
			p.repo.Commit()
		}
	}()

	if err = p.repo.InsertPayout(payoutModel); err != nil {
		return fmt.Errorf("Database payout insertion failed: %w", err)
	}

	for i := 1; i < len(transactions); i++ {
		donationModel := models.NewDonation(transactions[i].ID, 0, 0, 0, 0, "", "", sql.NullString{String: payoutModel.ID, Valid: true})
		updated, err := p.repo.UpdateRelatedPayout(donationModel)
		if err != nil {
			return fmt.Errorf("Update related payout failed for %s: %w", transactions[i].ID, err)
		}
		if updated {
			continue
		}

		charge, err := fetchRelatedCharge(transactions[i])
		if err != nil {
			return fmt.Errorf("Related charge fetch failed for %s: %w", transactions[i].ID, err)
		}

		if err = validateRelatedCharge(charge); err != nil {
			return fmt.Errorf("Related charge validation failed for %s: %w", transactions[i].ID, err)
		}

		donationModel = models.NewDonation(
			transactions[i].ID,
			uint64(transactions[i].Created),
			uint32(transactions[i].Amount),
			uint32(transactions[i].Fee),
			uint32(transactions[i].Net),
			charge.BillingDetails.Name,
			charge.BillingDetails.Email,
			sql.NullString{String: payoutModel.ID, Valid: true},
		)
		if err = p.repo.InsertDonation(donationModel); err != nil {
			return fmt.Errorf("Database donation insertion failed: %w", err)
		}
	}
	return nil
}

func validatePayout(payout *stripe.Payout) error {
	if payout.Status != "paid" {
		return fmt.Errorf(errors.ErrPayoutNotPaid)
	}
	if payout.ID == "" {
		return fmt.Errorf(errors.ErrPayoutIDMissing)
	}
	return nil
}

func fetchTransactions(id string) ([]*stripe.BalanceTransaction, error) {
	params := &stripe.BalanceTransactionListParams{}
	params.Payout = &id

	iter := balancetransaction.List(params)

	var transactions []*stripe.BalanceTransaction
	for iter.Next() {
		transactions = append(transactions, iter.BalanceTransaction())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}
	if len(transactions) < 2 {
		return nil, fmt.Errorf(errors.ErrTransactionListMissing)
	}
	return transactions, nil
}

func validatePayoutTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction.Type != "payout" {
		return fmt.Errorf(errors.ErrTransactionPayoutFailed+". Current type: %s", transaction.Type)
	}
	err := validateTransaction(transaction)
	return err
}

func validateChargeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction.Type != "charge" {
		return fmt.Errorf(errors.ErrTransactionChargeFailed+". Current type: %s", transaction.Type)
	}
	if transaction.Source == nil {
		return fmt.Errorf(errors.ErrChargeTransactionSourceMissing)
	}
	if transaction.Source.ID == "" {
		return fmt.Errorf(errors.ErrTransactionChargeIDMissing)
	}
	err := validateTransaction(transaction)
	return err
}

func validateMatchingSums(transactions []*stripe.BalanceTransaction) (int64, int64, int64, error) {
	var payoutGross, payoutFee int64

	for i := 1; i < len(transactions); i++ {
		payoutGross += transactions[i].Amount
		payoutFee += transactions[i].Fee
	}

	payoutNet := payoutGross - payoutFee
	payoutAmount := -transactions[0].Amount

	if payoutAmount != payoutNet {
		return 0, 0, 0, fmt.Errorf(errors.ErrTransactionPayoutMismatch)
	}
	return payoutGross, payoutFee, payoutNet, nil
}

func fetchRelatedCharge(transaction *stripe.BalanceTransaction) (*stripe.Charge, error) {
	params := &stripe.ChargeParams{}
	charge, err := charge.Get(transaction.Source.ID, params)
	if err != nil {
		return nil, err
	}
	return charge, nil
}

func validateRelatedCharge(charge *stripe.Charge) error {
	if charge.BillingDetails == nil {
		return fmt.Errorf(errors.ErrBillingDetailsMissing)
	}
	if charge.BillingDetails.Name == "" {
		return fmt.Errorf(errors.ErrClientNameMissing)
	}
	if charge.BillingDetails.Email == "" {
		return fmt.Errorf(errors.ErrClientEmailMissing)
	}
	return nil
}
