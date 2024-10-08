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
	GetRelatedDonations(payoutID string) ([]*models.Donation, error)
	GetPayout(id string) (*models.Payout, error)
	GetMonthlyPayouts(monthStart, monthEnd int64) ([]*models.Payout, error)
	GetRelatedFees(payoutID string) ([]*models.Fee, error)
}

type DocumentService interface {
	GenerateInvoice(donation *dto.FormattedDonation) (*gopdf.GoPdf, error)
	GeneratePayoutReport(payoutReportData *dto.PayoutReportData) (*gopdf.GoPdf, error)
	GenerateMonthlyReport(monthlyReportData *dto.MonthlyReportData) (*gopdf.GoPdf, error)
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

func (s *AccountingService) GenerateInvoice(id string) (pdf *gopdf.GoPdf, err error) {
	donationModel, err := s.repo.GetDonation(id)
	if err != nil {
		return nil, fmt.Errorf("fetch donation failed: %w", err)
	}
	donation := transformDonationModelToDTO(donationModel)
	pdf, err = s.document.GenerateInvoice(donation)
	if err != nil {
		return nil, fmt.Errorf("generate invoice failed: %w", err)
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
		return nil, fmt.Errorf("generate payout report failed: %w", err)
	}
	return
}

func (s *AccountingService) GenerateMonthlyReport(stringDate string) (pdf *gopdf.GoPdf, err error) {
	date, err := validateMonthString(stringDate)
	if err != nil {
		return nil, fmt.Errorf("month string invalid: %w", err)
	}

	monthStartUnix, monthEndUnix := getUnixTimestampsForMonth(date)
	payoutModels, err := s.repo.GetMonthlyPayouts(monthStartUnix, monthEndUnix)
	if err != nil {
		return nil, fmt.Errorf("fetch payouts failed: %w", err)
	}
	if len(payoutModels) == 0 {
		return nil, fmt.Errorf("monthly report empty")
	}

	gross, fee, net, err := monthlyReportSum(payoutModels)
	if err != nil {
		return nil, fmt.Errorf("monthly report sum failed: %w", err)
	}

	monthlyReportData := transformToMonthlyReportData(date, gross, fee, net, payoutModels)
	pdf, err = s.document.GenerateMonthlyReport(monthlyReportData)
	if err != nil {
		return nil, fmt.Errorf("generate monthly report failed: %w", err)
	}
	return
}

func (s *AccountingService) GenerateMonthlyReportView(stringDate string) (*dto.MonthlyReportView, error) {
	date, err := validateMonthString(stringDate)
	if err != nil {
		return nil, fmt.Errorf("month string invalid: %w", err)
	}

	monthStartUnix, monthEndUnix := getUnixTimestampsForMonth(date)
	payoutModels, err := s.repo.GetMonthlyPayouts(monthStartUnix, monthEndUnix)
	if err != nil {
		return nil, fmt.Errorf("fetch payouts failed: %w", err)
	}
	if len(payoutModels) == 0 {
		return transformToMonthlyReportView(stringDate, 0, 0, 0, nil), nil
	}

	gross, fee, net, err := monthlyReportSum(payoutModels)
	if err != nil {
		return nil, fmt.Errorf("monthly report sum failed: %w", err)
	}

	payouts, err := s.transformMonthlyViewPayoutModelsInDTOs(payoutModels)
	if err != nil {
		return nil, fmt.Errorf("monthly view payout models failed: %w", err)
	}

	return transformToMonthlyReportView(stringDate, gross, fee, net, payouts), nil
}

func (s AccountingService) transformMonthlyViewPayoutModelsInDTOs(payoutModels []*models.Payout) (payouts []*dto.FormattedPayout, err error) {
	for _, payoutModel := range payoutModels {
		donationModels, err := s.repo.GetRelatedDonations(payoutModel.ID)
		if err != nil {
			return nil, fmt.Errorf("fetch related donations failed: %w", err)
		}
		feeModels, err := s.repo.GetRelatedFees(payoutModel.ID)
		if err != nil {
			return nil, fmt.Errorf("fetch related fees failed: %w", err)
		}
		donations := transformDonationModelsToDTOs(donationModels)
		fees := transformFeeModelsToDTOs(feeModels)
		payouts = append(payouts, transformMonthlyViewPayoutModelToDTO(payoutModel, donations, fees))
	}
	return
}

func transformToMonthlyReportView(date string, gross, fee, net uint32, payouts []*dto.FormattedPayout) *dto.MonthlyReportView {
	return dto.NewMonthlyReportView(
		date,
		fmt.Sprintf("%.2f lei", float64(gross)/100),
		fmt.Sprintf("%.2f lei", float64(fee)/100),
		fmt.Sprintf("%.2f lei", float64(net)/100),
		payouts,
	)
}

func transformMonthlyViewPayoutModelToDTO(payout *models.Payout, donations []*dto.FormattedDonation, fees []*dto.FormattedFee) *dto.FormattedPayout {
	return dto.NewFormattedPayout(
		payout.ID,
		time.Unix(int64(payout.Created), 0).UTC().Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(payout.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Net)/100),
		donations,
		fees,
	)
}

func transformPayoutModelsToDTOs(payoutModels []*models.Payout) (payouts []*dto.FormattedPayout) {
	for _, payoutModel := range payoutModels {
		payouts = append(payouts, transformPayoutModelToDTO(payoutModel))
	}
	return
}

func transformPayoutModelToDTO(payout *models.Payout) *dto.FormattedPayout {
	return dto.NewFormattedPayout(
		payout.ID,
		time.Unix(int64(payout.Created), 0).UTC().Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(payout.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(payout.Net)/100),
		nil,
		nil,
	)
}

func transformDonationModelsToDTOs(donationModels []*models.Donation) (donations []*dto.FormattedDonation) {
	for _, donationModel := range donationModels {
		donations = append(donations, transformDonationModelToDTO(donationModel))
	}
	return
}

func transformDonationModelToDTO(donation *models.Donation) *dto.FormattedDonation {
	return dto.NewFormattedDonation(
		donation.ID,
		time.Unix(int64(donation.Created), 0).UTC().Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(donation.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Net)/100),
		donation.ClientName,
		donation.ClientEmail,
		donation.PayoutID.String,
	)
}

func transformDonationModelsToPayoutReportItems(donationModels []*models.Donation) (donations []*dto.PayoutReportItem) {
	for _, donationModel := range donationModels {
		donations = append(donations, transformDonationModelToPayoutReportItem(donationModel))
	}
	return
}

func transformDonationModelToPayoutReportItem(donation *models.Donation) *dto.PayoutReportItem {
	return dto.NewPayoutReportItem(
		donation.ID,
		"donation",
		"",
		time.Unix(int64(donation.Created), 0).UTC().Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(donation.Gross)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(donation.Net)/100),
	)
}

func transformFeeModelsToDTOs(feeModels []*models.Fee) (fees []*dto.FormattedFee) {
	for _, feeModel := range feeModels {
		fees = append(fees, transformFeeModelToDTO(feeModel))
	}
	return
}

func transformFeeModelToDTO(fee *models.Fee) *dto.FormattedFee {
	return dto.NewFormattedFee(
		fee.ID,
		fee.Description,
		time.Unix(int64(fee.Created), 0).UTC().Format("02 Jan 2006"),
		fmt.Sprintf("%.2f lei", float64(fee.Fee)/100),
	)
}

func transformFeeModelsToPayoutReportItems(feeModels []*models.Fee) (fees []*dto.PayoutReportItem) {
	for _, feeModel := range feeModels {
		fees = append(fees, transformFeeModelToPayoutReportItem(feeModel))
	}
	return
}

func transformFeeModelToPayoutReportItem(fee *models.Fee) *dto.PayoutReportItem {
	return dto.NewPayoutReportItem(
		fee.ID,
		"fee",
		fee.Description,
		time.Unix(int64(fee.Created), 0).UTC().Format("02 Jan 2006"),
		"0.00 lei",
		fmt.Sprintf("%.2f lei", float64(fee.Fee)/100),
		fmt.Sprintf("%.2f lei", float64(fee.Fee)/100),
	)
}

func transformToMonthlyReportData(date time.Time, gross, fee, net uint32, payoutModels []*models.Payout) *dto.MonthlyReportData {
	monthStart, monthEnd, emissionDate := getMonthDatesFromISO(date)
	payouts := transformPayoutModelsToDTOs(payoutModels)

	return dto.NewMonthlyReportData(
		monthStart,
		monthEnd,
		emissionDate,
		fmt.Sprintf("%.2f lei", float64(gross)/100),
		fmt.Sprintf("%.2f lei", float64(fee)/100),
		fmt.Sprintf("%.2f lei", float64(net)/100),
		payouts,
	)
}

func getUnixTimestampsForMonth(date time.Time) (monthStart, monthEnd int64) {
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthStart = start.Unix()

	end := start.AddDate(0, 1, 0).Add(-time.Second)
	monthEnd = end.Unix()
	return
}

func getMonthDatesFromISO(date time.Time) (monthStart, monthEnd, emissionDate string) {
	monthStartTime := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEndTime := monthStartTime.AddDate(0, 1, 0).Add(-time.Second)
	emissionDateTime := monthStartTime.AddDate(0, 1, 0)

	monthStart = monthStartTime.Format("2 Jan, 2006")
	monthEnd = monthEndTime.Format("2 Jan, 2006")
	emissionDate = emissionDateTime.Format("2 Jan, 2006")
	return
}

func validateMonthString(date string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}
	return parsedDate, nil
}

func monthlyReportSum(payouts []*models.Payout) (gross, fee, net uint32, err error) {
	for _, payout := range payouts {
		gross += payout.Gross
		fee += payout.Fee
		net += payout.Net
	}
	if gross-fee != net {
		return 0, 0, 0, fmt.Errorf("monthly gross-fee %v != net %v", gross-fee, net)
	}
	return
}
