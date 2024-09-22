package services

import (
	"database/sql"
	"fmt"

	"github.com/diother/go-invoices/internal/constants"
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
		return fmt.Errorf("payout validation error: %w", err)
	}

	transactions, err := fetchRelatedTransactions(payout.ID)
	if err != nil {
		return fmt.Errorf("related transactions fetch failed: %w", err)
	}

	if err = validateRelatedTransactions(transactions); err != nil {
		return fmt.Errorf("related transactions validation failed: %w", err)
	}

	payoutGross, payoutFee, payoutNet, err := validateMatchingSums(transactions)
	if err != nil {
		return fmt.Errorf("matching sum validation failed: %w", err)
	}

	if err = p.repo.BeginTransaction(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			p.repo.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
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
		return fmt.Errorf("database payout insertion failed: %w", err)
	}

	for _, transaction := range transactions[1:] {
		if err = p.PersistRelatedTransaction(transaction, payoutModel.ID); err != nil {
			return fmt.Errorf("related transaction persistence failed: %w", err)
		}
	}
	return
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

func (p *PayoutService) PersistRelatedTransaction(transaction *stripe.BalanceTransaction, payoutId string) (err error) {
	switch transaction.Type {
	case "charge":
		if err = p.UpsertDonation(transaction, payoutId); err != nil {
			return fmt.Errorf("upsert donation failed for %s: %w", transaction.ID, err)
		}

	case "stripe_fee":
		feeModel := models.NewFee(
			transaction.ID,
			uint64(transaction.Created),
			uint32(transaction.Amount),
			sql.NullString{String: payoutId, Valid: true},
		)
		if err = p.repo.InsertFee(feeModel); err != nil {
			return fmt.Errorf("database donation insertion failed: %w", err)
		}
	}
	return
}

func (p *PayoutService) UpsertDonation(transaction *stripe.BalanceTransaction, payoutId string) (err error) {
	donationModel := models.NewDonation(transaction.ID, 0, 0, 0, 0, "", "", sql.NullString{String: payoutId, Valid: true})
	updated, err := p.repo.UpdateRelatedPayout(donationModel)
	if err != nil {
		return fmt.Errorf("update related payout failed: %w", err)
	}
	if updated {
		return
	}

	charge, err := fetchRelatedCharge(transaction)
	if err != nil {
		return fmt.Errorf("related charge fetch failed: %w", err)
	}
	if err = validateCharge(charge); err != nil {
		return fmt.Errorf("related charge validation failed: %w", err)
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
		return fmt.Errorf("database donation insertion failed: %w", err)
	}
	return
}

func fetchRelatedCharge(transaction *stripe.BalanceTransaction) (*stripe.Charge, error) {
	params := &stripe.ChargeParams{}
	charge, err := charge.Get(transaction.Source.ID, params)
	if err != nil {
		return nil, err
	}
	return charge, nil
}

func validateRelatedTransactions(transactions []*stripe.BalanceTransaction) error {
	if len(transactions) < 2 {
		return fmt.Errorf(config.ErrPayoutListInsufficientTransactions)
	}
	if err := validatePayoutTransaction(transactions[0]); err != nil {
		return fmt.Errorf(config.ErrPayoutListPayoutTransactionInvalid+": %w", err)
	}
	for _, transaction := range transactions[1:] {
		if err := validateRelatedTransaction(transaction); err != nil {
			return fmt.Errorf(config.ErrPayoutListRelatedTransactionInvalid+": %w", err)
		}
	}
	return nil
}

func validateRelatedTransaction(transaction *stripe.BalanceTransaction) error {
	switch transaction.Type {
	case "charge":
		return validateChargeTransaction(transaction)
	case "stripe_fee":
		return validateFeeTransaction(transaction)
	default:
		return fmt.Errorf(config.ErrPayoutListUnexpectedTransaction+": %s", transaction.Type)
	}
}

func validateMatchingSums(transactions []*stripe.BalanceTransaction) (int64, int64, int64, error) {
	var payoutGross, payoutFee int64

	for _, transaction := range transactions[1:] {
		switch transaction.Type {
		case "charge":
			payoutGross += transaction.Amount
			payoutFee += transaction.Fee

		case "stripe_fee":
			payoutFee -= transaction.Amount
		}
	}

	payoutNet := payoutGross - payoutFee
	payoutAmount := -transactions[0].Amount

	if payoutAmount != payoutNet {
		return 0, 0, 0, fmt.Errorf(config.ErrPayoutListSumMismatch+". amount %v != net %v", payoutAmount, payoutNet)
	}
	return payoutGross, payoutFee, payoutNet, nil
}

func validatePayout(payout *stripe.Payout) error {
	if payout == nil {
		return fmt.Errorf(config.ErrPayoutMissing)
	}
	if payout.Status != "paid" {
		return fmt.Errorf(config.ErrPayoutStatusInvalid)
	}
	if payout.ID == "" {
		return fmt.Errorf(config.ErrPayoutIDMissing)
	}
	return nil
}

func validateCharge(charge *stripe.Charge) error {
	if charge == nil {
		return fmt.Errorf(config.ErrChargeMissing)
	}
	if charge.Status != "succeeded" {
		return fmt.Errorf(config.ErrChargeStatusInvalid)
	}
	if charge.BillingDetails == nil {
		return fmt.Errorf(config.ErrChargeBillingMissing)
	}
	if charge.BillingDetails.Name == "" {
		return fmt.Errorf(config.ErrChargeBillingNameMissing)
	}
	if charge.BillingDetails.Email == "" {
		return fmt.Errorf(config.ErrChargeBillingEmailMissing)
	}
	if charge.BalanceTransaction == nil {
		return fmt.Errorf(config.ErrTransactionMissing)
	}
	if charge.BalanceTransaction.ID == "" {
		return fmt.Errorf(config.ErrTransactionIDMissing)
	}
	return nil
}

func validatePayoutTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction == nil {
		return fmt.Errorf(config.ErrTransactionMissing)
	}
	if transaction.Type != "payout" {
		return fmt.Errorf(config.ErrPayoutTransactionTypeInvalid)
	}
	if transaction.ID == "" {
		return fmt.Errorf(config.ErrTransactionIDMissing)
	}
	if transaction.Created <= 0 {
		return fmt.Errorf(config.ErrTransactionCreatedInvalid)
	}
	if transaction.Amount >= 0 {
		return fmt.Errorf(config.ErrPayoutTransactionAmountInvalid)
	}
	if transaction.Fee != 0 {
		return fmt.Errorf(config.ErrPayoutTransactionFeeInvalid)
	}
	if transaction.Net >= 0 {
		return fmt.Errorf(config.ErrPayoutTransactionNetInvalid)
	}
	return nil
}

func validateChargeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction == nil {
		return fmt.Errorf(config.ErrTransactionMissing)
	}
	if transaction.Type != "charge" {
		return fmt.Errorf(config.ErrChargeTransactionTypeInvalid)
	}
	if transaction.ID == "" {
		return fmt.Errorf(config.ErrTransactionIDMissing)
	}
	if transaction.Created <= 0 {
		return fmt.Errorf(config.ErrTransactionCreatedInvalid)
	}
	if transaction.Amount <= 0 {
		return fmt.Errorf(config.ErrChargeTransactionAmountInvalid)
	}
	if transaction.Fee <= 0 {
		return fmt.Errorf(config.ErrChargeTransactionFeeInvalid)
	}
	if transaction.Net <= 0 {
		return fmt.Errorf(config.ErrChargeTransactionNetInvalid)
	}
	if transaction.Source == nil {
		return fmt.Errorf(config.ErrChargeTransactionSourceMissing)
	}
	if transaction.Source.ID == "" {
		return fmt.Errorf(config.ErrChargeTransactionSourceIDMissing)
	}
	return nil
}

func validateFeeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction == nil {
		return fmt.Errorf(config.ErrTransactionMissing)
	}
	if transaction.Type != "stripe_fee" {
		return fmt.Errorf(config.ErrFeeTransactionTypeInvalid)
	}
	if transaction.ID == "" {
		return fmt.Errorf(config.ErrTransactionIDMissing)
	}
	if transaction.Created <= 0 {
		return fmt.Errorf(config.ErrTransactionCreatedInvalid)
	}
	if transaction.Amount >= 0 {
		return fmt.Errorf(config.ErrFeeTransactionAmountInvalid)
	}
	if transaction.Fee != 0 {
		return fmt.Errorf(config.ErrFeeTransactionFeeInvalid)
	}
	if transaction.Net >= 0 {
		return fmt.Errorf(config.ErrFeeTransactionNetInvalid)
	}
	return nil
}
