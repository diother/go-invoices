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
	GetRelatedDonations(payoutID string) ([]*models.Donation, error)
	GetPayout(id string) (*models.Payout, error)
	GetAllPayouts() ([]*models.Payout, error)
	GetRelatedFees(payoutID string) ([]*models.Fee, error)
}

type DocumentService interface {
	GenerateInvoice(donation *dto.FormattedDonation) (*gopdf.GoPdf, error)
	GeneratePayoutReport(payoutReportData *dto.PayoutReportData) (*gopdf.GoPdf, error)
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
	donations = transformDonationModelsToDTOs(rawDonations)
	return
}

func (s *AccountingService) FetchPayouts() (payouts []*dto.FormattedPayout, err error) {
	rawPayouts, err := s.repo.GetAllPayouts()
	if err != nil {
		return nil, fmt.Errorf("fetch donations failed: %w", err)
	}
	payouts = transformPayoutModelsToDTOs(rawPayouts)
	return
}

func (s *AccountingService) GenerateInvoice(id string) (pdf *gopdf.GoPdf, err error) {
	donationModel, err := s.repo.GetDonation(id)
	if err != nil {
		return nil, fmt.Errorf("fetch donation failed: %w", err)
	}
	donation := transformDonationModelToDTO(donationModel)
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
	donationModels, err := s.repo.GetRelatedDonations(payoutID)
	if err != nil {
		return nil, fmt.Errorf("fetch related donations failed: %w", err)
	}
	feeModels, err := s.repo.GetRelatedFees(payoutID)
	if err != nil {
		return nil, fmt.Errorf("fetch related fees failed: %w", err)
	}

	items := transformDonationModelsToPayoutReportItems(donationModels)
	items = append(items, transformFeeModelsToPayoutReportItems(feeModels)...)

	payoutReportData := dto.NewPayoutReportData(
		transformPayoutModelToDTO(payoutModel),
		items,
	)
	pdf, err = s.document.GeneratePayoutReport(payoutReportData)
	if err != nil {
		return nil, fmt.Errorf("generating payout report failed: %w", err)
	}
	return
}

// needs unit test
func transformPayoutModelsToDTOs(payoutModels []*models.Payout) (payouts []*dto.FormattedPayout) {
	for _, payoutModel := range payoutModels {
		payouts = append(payouts, transformPayoutModelToDTO(payoutModel))
	}
	return
}

// needs unit test
func transformPayoutModelToDTO(payout *models.Payout) *dto.FormattedPayout {
	return dto.NewFormattedPayout(
		payout.ID,
		time.Unix(int64(payout.Created), 0).Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(payout.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Net)/100),
	)
}

// needs unit test
func transformDonationModelsToDTOs(donationModels []*models.Donation) (donations []*dto.FormattedDonation) {
	for _, donationModel := range donationModels {
		donations = append(donations, transformDonationModelToDTO(donationModel))
	}
	return
}

// needs unit test
func transformDonationModelToDTO(donation *models.Donation) *dto.FormattedDonation {
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

// needs unit test
func transformDonationModelsToPayoutReportItems(donationModels []*models.Donation) (donations []*dto.PayoutReportItem) {
	for _, donationModel := range donationModels {
		donations = append(donations, transformDonationModelToPayoutReportItem(donationModel))
	}
	return
}

// needs unit test
func transformDonationModelToPayoutReportItem(donation *models.Donation) *dto.PayoutReportItem {
	return dto.NewPayoutReportItem(
		donation.ID,
		"donation",
		time.Unix(int64(donation.Created), 0).Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(donation.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Net)/100),
	)
}

// needs unit test
func transformFeeModelsToDTOs(feeModels []*models.Fee) (fees []*dto.FormattedFee) {
	for _, feeModel := range feeModels {
		fees = append(fees, transformFeeModelToDTO(feeModel))
	}
	return
}

// needs unit test
func transformFeeModelToDTO(fee *models.Fee) *dto.FormattedFee {
	return dto.NewFormattedFee(
		fee.ID,
		time.Unix(int64(fee.Created), 0).Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(fee.Fee)/100),
	)
}

// needs unit test
func transformFeeModelsToPayoutReportItems(feeModels []*models.Fee) (fees []*dto.PayoutReportItem) {
	for _, feeModel := range feeModels {
		fees = append(fees, transformFeeModelToPayoutReportItem(feeModel))
	}
	return
}

// needs unit test
func transformFeeModelToPayoutReportItem(fee *models.Fee) *dto.PayoutReportItem {
	return dto.NewPayoutReportItem(
		fee.ID,
		"fee",
		time.Unix(int64(fee.Created), 0).Format("02 Jan 2006"),
		"0.00 lei",
		fmt.Sprintf("%.2f lei", float64(fee.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(fee.Fee)/100),
	)
}
