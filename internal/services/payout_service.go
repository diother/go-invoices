package services

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/diother/go-invoices/internal/models"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/balancetransaction"
	"github.com/stripe/stripe-go/v79/charge"
)

type PayoutRepository interface {
	Insert(payout *models.Payout) error
}

type PayoutServiceImpl struct {
	payoutRepo   PayoutRepository
	donationRepo DonationRepository
}

func NewPayoutServiceImpl(payoutRepo PayoutRepository, donationRepo DonationRepository) *PayoutServiceImpl {
	return &PayoutServiceImpl{
		payoutRepo:   payoutRepo,
		donationRepo: donationRepo,
	}
}

func (p *PayoutServiceImpl) ProcessPayout(payout *stripe.Payout) error {
	stripe.Key = "sk_test_51PVfUvDXCtuWOFq8ADmnd1iQEONLKIC6p1m1tALD67I6Ew4gRgOjoYGR7B5XK8hN0uc7iLE2Mbl9BedtgLIQubXU00XWzh1hmB"

	if payout.Status != "paid" {
		return errors.New("Payout was not paid")
	}

	payoutId := "po_1PvU4EDXCtuWOFq8p4uglXBe"
	if payoutId == "" {
		return errors.New("Payout.ID is missing")
	}

	transactions, err := fetchBalanceTransactions(payoutId)
	if err != nil {
		return err
	}

	var totalGross int64
	var totalFees int64

	for _, transaction := range transactions {
		if transaction.Type == "charge" {
			totalGross += transaction.Amount
			totalFees += transaction.Fee
		}
	}

	payoutTransaction := transactions[0]

	totalNet := totalGross - totalFees
	if -payoutTransaction.Amount != totalNet {
		return errors.New("Payout amount does not match total charges minus fees")
	}

	payoutModel := &models.Payout{
		ID:      payoutTransaction.ID,
		Created: uint64(payoutTransaction.Created),
		Gross:   uint32(totalGross),
		Fee:     uint32(totalFees),
		Net:     uint32(totalNet),
	}

	err = p.payoutRepo.Insert(payoutModel)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		if transaction.Type == "charge" {
			donationModel := &models.Donation{
				ID:       transaction.ID,
				PayoutID: payoutModel.ID,
			}
			updated, err := p.donationRepo.UpdatePayout(donationModel)
			if err != nil {
				return err
			}

			if !updated {
				params := &stripe.ChargeParams{}
				charge, err := charge.Get(transaction.Source.ID, params)
				if err != nil {
					return err
				}

				donationModel = &models.Donation{
					ID:          transaction.ID,
					Created:     uint64(transaction.Created),
					Gross:       uint32(transaction.Amount),
					Fee:         uint32(transaction.Fee),
					Net:         uint32(transaction.Net),
					ClientName:  charge.BillingDetails.Name,
					ClientEmail: charge.BillingDetails.Email,
					PayoutID:    payoutModel.ID,
				}

				err = p.donationRepo.Insert(donationModel)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *PayoutServiceImpl) buildPayoutModel(balanceTransaction *stripe.BalanceTransaction) *models.Payout {
	return &models.Payout{
		ID:      balanceTransaction.ID,
		Created: uint64(balanceTransaction.Created),
		Gross:   uint32(balanceTransaction.Amount),
		Fee:     uint32(balanceTransaction.Fee),
		Net:     uint32(balanceTransaction.Net),
	}
}

func fetchBalanceTransactions(id string) ([]*stripe.BalanceTransaction, error) {
	params := &stripe.BalanceTransactionListParams{
		Payout: &id,
	}
	i := balancetransaction.List(params)

	var transactions []*stripe.BalanceTransaction
	for i.Next() {
		transactions = append(transactions, i.BalanceTransaction())
	}

	if len(transactions) == 0 {
		return nil, errors.New("No transactions were fetched")
	}
	return transactions, nil
}

type Payout struct {
	payout *stripe.Payout
}

func NewPayout(payout *stripe.Payout) *Payout {
	return &Payout{payout: payout}
}

type Transaction struct {
	transaction *stripe.BalanceTransaction
}

func NewTransaction(transaction *stripe.BalanceTransaction) *Transaction {
	return &Transaction{transaction: transaction}
}

func (p *Payout) Log() {
	responseJSON, err := json.MarshalIndent(p.payout, "", "    ")
	if err != nil {
		log.Printf("Error marshalling payout: %v\n", err)
	}
	log.Println("Payout details:", string(responseJSON))
}

func (t *Transaction) Log() {
	responseJSON, err := json.MarshalIndent(t.transaction, "", "    ")
	if err != nil {
		log.Printf("Error marshalling payout: %v\n", err)
	}
	log.Println("Type:", t.transaction.Type)
	log.Println("Transaction details:", string(responseJSON))
}
