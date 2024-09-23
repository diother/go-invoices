package services

import (
	"fmt"

	"github.com/diother/go-invoices/internal/models"
	"github.com/signintech/gopdf"
)

type AccountingService struct {
	repo WebhookRepository
}

func NewAccountingService(repo WebhookRepository) *AccountingService {
	return &AccountingService{repo: repo}
}

func (s *AccountingService) FetchDonations() (donations []*models.Donation, err error) {
	donations, err = s.repo.GetAllDonations()
	if err != nil {
		return nil, fmt.Errorf("Fetch donations failed: %w", err)
	}
	return donations, nil
}

func (s *AccountingService) GenerateDocument(documentType, documentID string) (pdf gopdf.GoPdf, err error) {
	pdf = gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err = pdf.AddTTFFont("Roboto", "./internal/pdf/static/fonts/Roboto-Regular.ttf"); err != nil {
		return
	}

	fmt.Println("Type:", documentType)
	fmt.Println("ID:", documentID)

	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(94, 100, 112)

	pdf.SetXY(25, 25)
	pdf.Cell(nil, "Hello, World!")
	return
}
