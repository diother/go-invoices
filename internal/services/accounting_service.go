package services

import (
	"fmt"
	"time"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/diother/go-invoices/internal/models"
	"github.com/signintech/gopdf"
)

type AccountingService struct {
	repo WebhookRepository
}

func NewAccountingService(repo WebhookRepository) *AccountingService {
	return &AccountingService{repo: repo}
}

func (s *AccountingService) FetchDonations() (donations []*dto.FormattedDonation, err error) {
	rawDonations, err := s.repo.GetAllDonations()
	if err != nil {
		return nil, fmt.Errorf("Fetch donations failed: %w", err)
	}
	for _, rawDonation := range rawDonations {
		donations = append(donations, formatDonationModel(rawDonation))
	}
	return donations, nil
}

func (s *AccountingService) GenerateInvoice(id string) (pdf *gopdf.GoPdf, err error) {
	donationModel, err := s.repo.GetDonation(id)
	if err != nil {
		return nil, fmt.Errorf("Fetch donations failed: %w", err)
	}

	donation := formatDonationModel(donationModel)

	fmt.Println(donation)

	pdf = &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err = pdf.AddTTFFont("Roboto", "./internal/pdf/static/fonts/Roboto-Regular.ttf"); err != nil {
		return
	}

	fmt.Println("ID:", id)

	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(94, 100, 112)

	pdf.SetXY(25, 25)
	pdf.Cell(nil, "Hello, World!")
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
