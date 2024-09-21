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
	InsertFee(fee *models.Fee) error
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

func (p *PayoutService) ProcessPayout(payout *stripe.Payout) (err error) {
	if err = validatePayout(payout); err != nil {
		return fmt.Errorf("Payout validation error: %w", err)
	}

	transactions, err := fetchRelatedTransactions(payout.ID)
	if err != nil {
		return fmt.Errorf("Related transactions fetch failed: %w", err)
	}

	if err = validateRelatedTransactions(transactions); err != nil {
		return fmt.Errorf("Related transactions validation failed: %w", err)
	}

	payoutGross, payoutFee, payoutNet, err := validateMatchingSums(transactions)
	if err != nil {
		return fmt.Errorf("Matching sum validation failed: %w", err)
	}

	if err = p.repo.BeginTransaction(); err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			p.repo.Rollback()
			err = fmt.Errorf("Panic occurred: %v", r)
		} else if err != nil {
			p.repo.Rollback()
		} else {
			p.repo.Commit()
		}
	}()

	payoutModel := models.NewPayout(
		transactions[0].ID,
		uint64(transactions[0].Created),
		uint32(payoutGross),
		uint32(payoutFee),
		uint32(payoutNet),
	)
	if err = p.repo.InsertPayout(payoutModel); err != nil {
		return fmt.Errorf("Database payout insertion failed: %w", err)
	}

	for i := 1; i < len(transactions); i++ {
		switch transactions[i].Type {
		case "charge":
			if err = p.UpsertDonation(transactions[i], payoutModel.ID); err != nil {
				return fmt.Errorf("Upsert donation failed for %s: %w", transactions[i].ID, err)
			}

		case "stripe_fee":
			feeModel := models.NewFee(
				transactions[i].ID,
				uint64(transactions[i].Created),
				uint32(transactions[i].Amount),
				sql.NullString{String: payoutModel.ID, Valid: true},
			)
			if err = p.repo.InsertFee(feeModel); err != nil {
				return fmt.Errorf("Database donation insertion failed: %w", err)
			}
		}
	}
	return
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

func fetchRelatedTransactions(id string) ([]*stripe.BalanceTransaction, error) {
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
	return transactions, nil
}

func validateRelatedTransactions(transactions []*stripe.BalanceTransaction) error {
	if len(transactions) < 2 {
		return fmt.Errorf(errors.ErrTransactionListMissing)
	}
	if err := validatePayoutTransaction(transactions[0]); err != nil {
		return fmt.Errorf("Payout transaction validation failed: %w", err)
	}
	for i := 1; i < len(transactions); i++ {
		switch transactions[i].Type {
		case "charge":
			if err := validateChargeTransaction(transactions[i]); err != nil {
				return fmt.Errorf("Charge transaction validation failed for %s: %w", transactions[i].ID, err)
			}
		case "stripe_fee":
			if err := validateFeeTransaction(transactions[i]); err != nil {
				return fmt.Errorf("Fee transaction validation failed for %s: %w", transactions[i].ID, err)
			}
		default:
			return fmt.Errorf("Unexpected transaction type for %s: %s", transactions[i].ID, transactions[i].Type)
		}
	}
	return nil
}

func validatePayoutTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction.Type != "payout" {
		return fmt.Errorf(errors.ErrTransactionPayoutFailed+". Current type: %s", transaction.Type)
	}
	return validateTransaction(transaction)
}

func validateChargeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction.Source == nil {
		return fmt.Errorf(errors.ErrChargeTransactionSourceMissing)
	}
	if transaction.Source.ID == "" {
		return fmt.Errorf(errors.ErrTransactionChargeIDMissing)
	}
	return validateTransaction(transaction)
}

func validateFeeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction.ID == "" {
		return fmt.Errorf(errors.ErrTransactionIDMissing)
	}
	if transaction.Created == 0 {
		return fmt.Errorf(errors.ErrTransactionCreatedMissing)
	}
	if transaction.Amount == 0 {
		return fmt.Errorf(errors.ErrTransactionAmountMissing)
	}
	return nil
}

func validateMatchingSums(transactions []*stripe.BalanceTransaction) (int64, int64, int64, error) {
	var payoutGross, payoutFee int64
	for i := 1; i < len(transactions); i++ {
		switch transactions[i].Type {
		case "charge":
			payoutGross += transactions[i].Amount
			payoutFee += transactions[i].Fee

		case "stripe_fee":
			payoutFee -= transactions[i].Amount
		}
	}

	payoutNet := payoutGross - payoutFee
	payoutAmount := -transactions[0].Amount

	if payoutAmount != payoutNet {
		return 0, 0, 0, fmt.Errorf(errors.ErrTransactionPayoutMismatch)
	}
	return payoutGross, payoutFee, payoutNet, nil
}

func (p *PayoutService) UpsertDonation(transaction *stripe.BalanceTransaction, payoutId string) (err error) {
	donationModel := models.NewDonation(transaction.ID, 0, 0, 0, 0, "", "", sql.NullString{String: payoutId, Valid: true})
	updated, err := p.repo.UpdateRelatedPayout(donationModel)
	if err != nil {
		return fmt.Errorf("Update related payout failed: %w", err)
	}
	if updated {
		return
	}

	charge, err := fetchRelatedCharge(transaction)
	if err != nil {
		return fmt.Errorf("Related charge fetch failed: %w", err)
	}
	if err = validateRelatedCharge(charge); err != nil {
		return fmt.Errorf("Related charge validation failed: %w", err)
	}
	donationModel = models.NewDonation(
		transaction.ID,
		uint64(transaction.Created),
		uint32(transaction.Amount),
		uint32(transaction.Fee),
		uint32(transaction.Net),
		charge.BillingDetails.Name,
		charge.BillingDetails.Email,
		sql.NullString{String: payoutId, Valid: true},
	)
	if err = p.repo.InsertDonation(donationModel); err != nil {
		return fmt.Errorf("Database donation insertion failed: %w", err)
	}
	return nil
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
