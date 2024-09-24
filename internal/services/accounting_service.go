package services

import (
	"fmt"
	"time"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/diother/go-invoices/internal/models"
	"github.com/signintech/gopdf"
)

type PWARepository interface {
	GetDonation(id string) (*models.Donation, error)
	GetAllDonations() ([]*models.Donation, error)
	// GetRelatedDonations(payoutID string) ([]*models.Donation, error)
	GetPayout(id string) (*models.Payout, error)
	GetAllPayouts() ([]*models.Payout, error)
	// GetRelatedFees(payoutID string) ([]*models.Donation, error)
}

type DocumentService interface {
	GenerateInvoice(donation *dto.FormattedDonation) (*gopdf.GoPdf, error)
}

type AccountingService struct {
	repo     PWARepository
	document DocumentService
}

func NewAccountingService(repo PWARepository, document DocumentService) *AccountingService {
	return &AccountingService{
		repo:     repo,
		document: document,
	}
}

func (s *AccountingService) FetchDonations() (donations []*dto.FormattedDonation, err error) {
	rawDonations, err := s.repo.GetAllDonations()
	if err != nil {
		return nil, fmt.Errorf("fetch donations failed: %w", err)
	}
	donations = formatDonations(rawDonations)
	return
}

func (s *AccountingService) FetchPayouts() (payouts []*dto.FormattedPayout, err error) {
	rawPayouts, err := s.repo.GetAllPayouts()
	if err != nil {
		return nil, fmt.Errorf("fetch donations failed: %w", err)
	}
	payouts = formatPayouts(rawPayouts)
	return
}

func (s *AccountingService) GenerateInvoice(id string) (pdf *gopdf.GoPdf, err error) {
	donationModel, err := s.repo.GetDonation(id)
	if err != nil {
		return nil, fmt.Errorf("fetch donation failed: %w", err)
	}
	donation := formatDonationModel(donationModel)
	pdf, err = s.document.GenerateInvoice(donation)
	if err != nil {
		return nil, fmt.Errorf("generating invoice failed: %w", err)
	}
	return
}

func (s *AccountingService) GeneratePayoutReport(payoutID string) (pdf *gopdf.GoPdf, err error) {
	payoutModel, err := s.repo.GetPayout(payoutID)
	if err != nil {
		return nil, fmt.Errorf("fetch payout failed: %w", err)
	}
	// donations, err := s.repo.GetRelatedDonations(payoutID)
	// if err != nil {
	// 	return nil, fmt.Errorf("fetch related donations failed: %w", err)
	// }
	fmt.Println(payoutModel)
	// fees, err := s.repo.GetRelatedFees(payoutID)
	// query the donations associated with the payout
	// query the stripe_fees associated with the payout
	// format the payout
	// format the donations
	// format the stripe_fees
	// generate the payout report
	return
}

func formatDonations(rawDonations []*models.Donation) (donations []*dto.FormattedDonation) {
	for _, rawDonation := range rawDonations {
		donations = append(donations, formatDonationModel(rawDonation))
	}
	return
}

func formatDonationModel(donation *models.Donation) *dto.FormattedDonation {
	return dto.NewFormattedDonation(
		donation.ID,
		time.Unix(int64(donation.Created), 0).Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(donation.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Net)/100),
		donation.ClientName,
		donation.ClientEmail,
		donation.PayoutID.String,
	)
}

func formatPayouts(rawPayouts []*models.Payout) (payouts []*dto.FormattedPayout) {
	for _, rawPayout := range rawPayouts {
		payouts = append(payouts, formatPayoutModel(rawPayout))
	}
	return
}

func formatPayoutModel(payout *models.Payout) *dto.FormattedPayout {
	return dto.NewFormattedPayout(
		payout.ID,
		time.Unix(int64(payout.Created), 0).Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(payout.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Net)/100),
	)
}
