package services

import (
	"errors"

	"github.com/diother/go-invoices/internal/models"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
)

type DonationRepository interface {
	Insert(donation *models.Donation) error
	UpdatePayout(donation *models.Donation) (bool, error)
}

type DonationServiceImpl struct {
	repo DonationRepository
}

func NewDonationServiceImpl(repo DonationRepository) *DonationServiceImpl {
	return &DonationServiceImpl{repo: repo}
}

func (d *DonationServiceImpl) ProcessDonation(charge *stripe.Charge) error {
	stripe.Key = "sk_test_51PVfUvDXCtuWOFq8ADmnd1iQEONLKIC6p1m1tALD67I6Ew4gRgOjoYGR7B5XK8hN0uc7iLE2Mbl9BedtgLIQubXU00XWzh1hmB"

	if charge.Status != "succeeded" {
		return errors.New("Charge did not succeed")
	}

	clientName := charge.BillingDetails.Name
	clientEmail := charge.BillingDetails.Email
	if clientName == "" || clientEmail == "" {
		return errors.New("Client name or email is missing")
	}

	paymentIntentId := charge.PaymentIntent.ID
	balanceTransactionId := charge.BalanceTransaction.ID
	if paymentIntentId == "" || balanceTransactionId == "" {
		return errors.New("PaymentIntent.ID or BalanceTransaction.ID is missing")
	}

	balanceTransaction, err := fetchBalanceTransaction(balanceTransactionId)
	if err != nil {
		return err
	}

	donationModel := d.buildDonationModel(charge, balanceTransaction)
	err = d.repo.Insert(donationModel)
	if err != nil {
		return err
	}

	return nil
}

func fetchBalanceTransaction(id string) (*stripe.BalanceTransaction, error) {
	params := &stripe.BalanceTransactionParams{}
	balanceTransaction, err := balancetransaction.Get(id, params)
	if err != nil {
		return nil, err
	}
	return balanceTransaction, nil
}

func (d *DonationServiceImpl) buildDonationModel(charge *stripe.Charge, balanceTransaction *stripe.BalanceTransaction) *models.Donation {
	return &models.Donation{
		ID:          balanceTransaction.ID,
		Created:     uint64(balanceTransaction.Created),
		Gross:       uint32(balanceTransaction.Amount),
		Fee:         uint32(balanceTransaction.Fee),
		Net:         uint32(balanceTransaction.Net),
		ClientName:  charge.BillingDetails.Name,
		ClientEmail: charge.BillingDetails.Email,
	}
}
