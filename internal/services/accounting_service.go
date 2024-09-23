package services

import (
	"fmt"
	"time"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/diother/go-invoices/internal/models"
	"github.com/signintech/gopdf"
)

type DocumentService interface {
	GenerateInvoice(donation *dto.FormattedDonation) (*gopdf.GoPdf, error)
}

type AccountingService struct {
	repo     WebhookRepository
	document DocumentService
}

func NewAccountingService(repo WebhookRepository, document DocumentService) *AccountingService {
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
	for _, rawDonation := range rawDonations {
		donations = append(donations, formatDonationModel(rawDonation))
	}
	return donations, nil
}

func (s *AccountingService) GenerateInvoice(id string) (pdf *gopdf.GoPdf, err error) {
	donationModel, err := s.repo.GetDonation(id)
	if err != nil {
		return nil, fmt.Errorf("fetch donations failed: %w", err)
	}

	donation := formatDonationModel(donationModel)

	pdf, err = s.document.GenerateInvoice(donation)
	if err != nil {
		return nil, fmt.Errorf("generating invoice failed: %w", err)
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
