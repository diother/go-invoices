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

func (s *PayoutService) ProcessPayout(payout *stripe.Payout) (err error) {
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

	if err = s.repo.BeginTransaction(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			s.repo.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			s.repo.Rollback()
		} else {
			s.repo.Commit()
		}
	}()

	payoutModel := transformPayoutDTOToModel(transactions[0], payoutGross, payoutFee, payoutNet)
	if err = s.repo.InsertPayout(payoutModel); err != nil {
		return fmt.Errorf("database payout insertion failed: %w", err)
	}

	for _, transaction := range transactions[1:] {
		if err = s.PersistRelatedTransaction(transaction, payoutModel.ID); err != nil {
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

func (s *PayoutService) PersistRelatedTransaction(transaction *stripe.BalanceTransaction, payoutID string) (err error) {
	switch transaction.Type {
	case "charge":
		if err = s.UpsertDonation(transaction, payoutID); err != nil {
			return fmt.Errorf("upsert donation failed for %s: %w", transaction.ID, err)
		}

	case "stripe_fee":
		feeModel := transformFeeDTOToModel(transaction, payoutID)
		if err = s.repo.InsertFee(feeModel); err != nil {
			return fmt.Errorf("database donation insertion failed: %w", err)
		}
	}
	return
}

func (s *PayoutService) UpsertDonation(transaction *stripe.BalanceTransaction, payoutID string) (err error) {
	donationModel := transformUpdateDonationDTOToModel(transaction.ID, payoutID)
	updated, err := s.repo.UpdateRelatedPayout(donationModel)
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

	donationModel = transformDonationDTOToModel(transaction, charge, payoutID)
	if err = s.repo.InsertDonation(donationModel); err != nil {
		return fmt.Errorf("database donation insertion failed: %w", err)
	}
	return
}

func transformPayoutDTOToModel(transaction *stripe.BalanceTransaction, gross, fee, net int64) *models.Payout {
	return models.NewPayout(
		transaction.ID,
		uint64(transaction.Created),
		uint32(gross),
		uint32(fee),
		uint32(net),
	)
}

func transformUpdateDonationDTOToModel(transactionID, payoutID string) *models.Donation {
	return models.NewDonation(transactionID, 0, 0, 0, 0, "", "", sql.NullString{String: payoutID, Valid: true})
}

func transformDonationDTOToModel(transaction *stripe.BalanceTransaction, charge *stripe.Charge, payoutID string) *models.Donation {
	return models.NewDonation(
		transaction.ID,
		uint64(transaction.Created),
		uint32(transaction.Amount),
		uint32(transaction.Fee),
		uint32(transaction.Net),
		charge.BillingDetails.Name,
		charge.BillingDetails.Email,
		sql.NullString{String: payoutID, Valid: true},
	)
}

func transformFeeDTOToModel(transaction *stripe.BalanceTransaction, payoutID string) *models.Fee {
	return models.NewFee(
		transaction.ID,
		transaction.Description,
		uint64(transaction.Created),
		uint32(-transaction.Amount),
		sql.NullString{String: payoutID, Valid: true},
	)
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
		return fmt.Errorf(constants.ErrPayoutListInsufficientTransactions)
	}
	if err := validatePayoutTransaction(transactions[0]); err != nil {
		return fmt.Errorf(constants.ErrPayoutListPayoutTransactionInvalid+": %w", err)
	}
	for _, transaction := range transactions[1:] {
		if err := validateRelatedTransaction(transaction); err != nil {
			return fmt.Errorf(constants.ErrPayoutListRelatedTransactionInvalid+": %w", err)
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
		return fmt.Errorf(constants.ErrPayoutListUnexpectedTransaction+": %s", transaction.Type)
	}
}

func validateMatchingSums(transactions []*stripe.BalanceTransaction) (payoutGross, payoutFee, payoutNet int64, err error) {
	for _, transaction := range transactions[1:] {
		switch transaction.Type {
		case "charge":
			payoutGross += transaction.Amount
			payoutFee += transaction.Fee

		case "stripe_fee":
			payoutFee -= transaction.Amount
		}
	}
	payoutNet = payoutGross - payoutFee
	payoutAmount := -transactions[0].Amount

	if payoutAmount != payoutNet {
		return 0, 0, 0, fmt.Errorf(constants.ErrPayoutListSumMismatch+". amount %v != net %v", payoutAmount, payoutNet)
	}
	return
}

func validatePayout(payout *stripe.Payout) error {
	if payout == nil {
		return fmt.Errorf(constants.ErrPayoutMissing)
	}
	if payout.Status != "paid" {
		return fmt.Errorf(constants.ErrPayoutStatusInvalid)
	}
	if payout.ID == "" {
		return fmt.Errorf(constants.ErrPayoutIDMissing)
	}
	return nil
}

func validateCharge(charge *stripe.Charge) error {
	if charge == nil {
		return fmt.Errorf(constants.ErrChargeMissing)
	}
	if charge.Status != "succeeded" {
		return fmt.Errorf(constants.ErrChargeStatusInvalid)
	}
	if charge.BillingDetails == nil {
		return fmt.Errorf(constants.ErrChargeBillingMissing)
	}
	if charge.BillingDetails.Name == "" {
		return fmt.Errorf(constants.ErrChargeBillingNameMissing)
	}
	if charge.BillingDetails.Email == "" {
		return fmt.Errorf(constants.ErrChargeBillingEmailMissing)
	}
	if charge.BalanceTransaction == nil {
		return fmt.Errorf(constants.ErrTransactionMissing)
	}
	if charge.BalanceTransaction.ID == "" {
		return fmt.Errorf(constants.ErrTransactionIDMissing)
	}
	return nil
}

func validatePayoutTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction == nil {
		return fmt.Errorf(constants.ErrTransactionMissing)
	}
	if transaction.Type != "payout" {
		return fmt.Errorf(constants.ErrPayoutTransactionTypeInvalid)
	}
	if transaction.ID == "" {
		return fmt.Errorf(constants.ErrTransactionIDMissing)
	}
	if transaction.Created <= 0 {
		return fmt.Errorf(constants.ErrTransactionCreatedInvalid)
	}
	if transaction.Amount >= 0 {
		return fmt.Errorf(constants.ErrPayoutTransactionAmountInvalid)
	}
	if transaction.Fee != 0 {
		return fmt.Errorf(constants.ErrPayoutTransactionFeeInvalid)
	}
	if transaction.Net >= 0 {
		return fmt.Errorf(constants.ErrPayoutTransactionNetInvalid)
	}
	return nil
}

func validateChargeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction == nil {
		return fmt.Errorf(constants.ErrTransactionMissing)
	}
	if transaction.Type != "charge" {
		return fmt.Errorf(constants.ErrChargeTransactionTypeInvalid)
	}
	if transaction.ID == "" {
		return fmt.Errorf(constants.ErrTransactionIDMissing)
	}
	if transaction.Created <= 0 {
		return fmt.Errorf(constants.ErrTransactionCreatedInvalid)
	}
	if transaction.Amount <= 0 {
		return fmt.Errorf(constants.ErrChargeTransactionAmountInvalid)
	}
	if transaction.Fee <= 0 {
		return fmt.Errorf(constants.ErrChargeTransactionFeeInvalid)
	}
	if transaction.Net <= 0 {
		return fmt.Errorf(constants.ErrChargeTransactionNetInvalid)
	}
	if transaction.Source == nil {
		return fmt.Errorf(constants.ErrChargeTransactionSourceMissing)
	}
	if transaction.Source.ID == "" {
		return fmt.Errorf(constants.ErrChargeTransactionSourceIDMissing)
	}
	return nil
}

func validateFeeTransaction(transaction *stripe.BalanceTransaction) error {
	if transaction == nil {
		return fmt.Errorf(constants.ErrTransactionMissing)
	}
	if transaction.Type != "stripe_fee" {
		return fmt.Errorf(constants.ErrFeeTransactionTypeInvalid)
	}
	if transaction.ID == "" {
		return fmt.Errorf(constants.ErrTransactionIDMissing)
	}
	if transaction.Description == "" {
		return fmt.Errorf(constants.ErrFeeTransactionDescriptionMissing)
	}
	if transaction.Created <= 0 {
		return fmt.Errorf(constants.ErrTransactionCreatedInvalid)
	}
	if transaction.Amount >= 0 {
		return fmt.Errorf(constants.ErrFeeTransactionAmountInvalid)
	}
	if transaction.Fee != 0 {
		return fmt.Errorf(constants.ErrFeeTransactionFeeInvalid)
	}
	if transaction.Net >= 0 {
		return fmt.Errorf(constants.ErrFeeTransactionNetInvalid)
	}
	return nil
}
