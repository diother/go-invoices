package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/diother/go-invoices/internal/models"
)

func TestTransformToMonthlyReportView(t *testing.T) {
	testCases := map[string]struct {
		date     string
		gross    uint32
		fee      uint32
		net      uint32
		payouts  []*dto.FormattedPayout
		expected *dto.MonthlyReportView
	}{
		"validReport": {
			date:    "2024-09",
			gross:   123456,
			fee:     1234,
			net:     122222,
			payouts: []*dto.FormattedPayout{{ID: "payout1"}},
			expected: &dto.MonthlyReportView{
				Date:  "2024-09",
				Gross: "1234.56 lei",
				Fee:   "12.34 lei",
				Net:   "1222.22 lei",
				Payouts: []*dto.FormattedPayout{
					{ID: "payout1"},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformToMonthlyReportView(tc.date, tc.gross, tc.fee, tc.net, tc.payouts)
			if result.Date != tc.expected.Date || result.Gross != tc.expected.Gross ||
				result.Fee != tc.expected.Fee || result.Net != tc.expected.Net {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
			if len(result.Payouts) != len(tc.expected.Payouts) {
				t.Errorf("Expected %d payouts, got %d", len(tc.expected.Payouts), len(result.Payouts))
			}
		})
	}
}

func TestTransformMonthlyViewPayoutModelToDTO(t *testing.T) {
	testCases := map[string]struct {
		payout    *models.Payout
		donations []*dto.FormattedDonation
		fees      []*dto.FormattedFee
		expected  *dto.FormattedPayout
	}{
		"validPayoutDTO": {
			payout: &models.Payout{
				ID:      "payout1",
				Created: 1700000000,
				Gross:   10000,
				Fee:     100,
				Net:     9900,
			},
			donations: []*dto.FormattedDonation{{ID: "donation1"}},
			fees:      []*dto.FormattedFee{{ID: "fee1"}},
			expected: &dto.FormattedPayout{
				ID:      "payout1",
				Created: "14 Nov 2023",
				Gross:   "100.00 lei",
				Fee:     "1.00 lei",
				Net:     "99.00 lei",
				Donations: []*dto.FormattedDonation{
					{ID: "donation1"},
				},
				Fees: []*dto.FormattedFee{
					{ID: "fee1"},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformMonthlyViewPayoutModelToDTO(tc.payout, tc.donations, tc.fees)

			if result.ID != tc.expected.ID || result.Created != tc.expected.Created ||
				result.Gross != tc.expected.Gross || result.Fee != tc.expected.Fee ||
				result.Net != tc.expected.Net {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
			if len(result.Donations) != len(tc.expected.Donations) {
				t.Errorf("Expected %d donations, got %d", len(tc.expected.Donations), len(result.Donations))
			}
			if len(result.Fees) != len(tc.expected.Fees) {
				t.Errorf("Expected %d fees, got %d", len(tc.expected.Fees), len(result.Fees))
			}
		})
	}
}

func TestTransformPayoutModelToDTO(t *testing.T) {
	testCases := map[string]struct {
		input    *models.Payout
		expected *dto.FormattedPayout
	}{
		"validPayout": {
			input: &models.Payout{
				ID:      "payout1",
				Created: 1700000000,
				Gross:   15000,
				Fee:     500,
				Net:     14500,
			},
			expected: &dto.FormattedPayout{
				ID:        "payout1",
				Created:   "14 Nov 2023",
				Gross:     "150.00 lei",
				Fee:       "5.00 lei",
				Net:       "145.00 lei",
				Donations: nil,
				Fees:      nil,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformPayoutModelToDTO(tc.input)

			if result.ID != tc.expected.ID ||
				result.Created != tc.expected.Created ||
				result.Gross != tc.expected.Gross ||
				result.Fee != tc.expected.Fee ||
				result.Net != tc.expected.Net {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTransformPayoutModelsToDTOs(t *testing.T) {
	testCases := map[string]struct {
		input    []*models.Payout
		expected []*dto.FormattedPayout
	}{
		"multiplePayouts": {
			input: []*models.Payout{
				{
					ID:      "payout1",
					Created: 1700000000,
					Gross:   15000,
					Fee:     500,
					Net:     14500,
				},
				{
					ID:      "payout2",
					Created: 1700000500,
					Gross:   20000,
					Fee:     600,
					Net:     19400,
				},
			},
			expected: []*dto.FormattedPayout{
				{
					ID:        "payout1",
					Created:   "14 Nov 2023",
					Gross:     "150.00 lei",
					Fee:       "5.00 lei",
					Net:       "145.00 lei",
					Donations: nil,
					Fees:      nil,
				},
				{
					ID:        "payout2",
					Created:   "14 Nov 2023",
					Gross:     "200.00 lei",
					Fee:       "6.00 lei",
					Net:       "194.00 lei",
					Donations: nil,
					Fees:      nil,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformPayoutModelsToDTOs(tc.input)

			if len(result) != len(tc.expected) {
				t.Fatalf("Expected %d results, got %d", len(tc.expected), len(result))
			}

			for i, res := range result {
				if res.ID != tc.expected[i].ID ||
					res.Created != tc.expected[i].Created ||
					res.Gross != tc.expected[i].Gross ||
					res.Fee != tc.expected[i].Fee ||
					res.Net != tc.expected[i].Net {
					t.Errorf("Expected %v, got %v", tc.expected[i], res)
				}
			}
		})
	}
}

func TestTransformDonationModelToDTO(t *testing.T) {
	testCases := map[string]struct {
		input    *models.Donation
		expected *dto.FormattedDonation
	}{
		"validDonation": {
			input: &models.Donation{
				ID:          "donation1",
				Created:     1700000000,
				Gross:       5000,
				Fee:         100,
				Net:         4900,
				ClientName:  "John Doe",
				ClientEmail: "john@example.com",
				PayoutID:    sql.NullString{String: "payout1", Valid: true},
			},
			expected: &dto.FormattedDonation{
				ID:          "donation1",
				Created:     "14 Nov 2023",
				Gross:       "50.00 lei",
				Fee:         "1.00 lei",
				Net:         "49.00 lei",
				ClientName:  "John Doe",
				ClientEmail: "john@example.com",
				PayoutID:    "payout1",
			},
		},
		"donationWithoutPayoutID": {
			input: &models.Donation{
				ID:          "donation2",
				Created:     1700000500,
				Gross:       10000,
				Fee:         500,
				Net:         9500,
				ClientName:  "Jane Doe",
				ClientEmail: "jane@example.com",
				PayoutID:    sql.NullString{Valid: false},
			},
			expected: &dto.FormattedDonation{
				ID:          "donation2",
				Created:     "14 Nov 2023",
				Gross:       "100.00 lei",
				Fee:         "5.00 lei",
				Net:         "95.00 lei",
				ClientName:  "Jane Doe",
				ClientEmail: "jane@example.com",
				PayoutID:    "",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformDonationModelToDTO(tc.input)

			if result.ID != tc.expected.ID ||
				result.Created != tc.expected.Created ||
				result.Gross != tc.expected.Gross ||
				result.Fee != tc.expected.Fee ||
				result.Net != tc.expected.Net ||
				result.ClientName != tc.expected.ClientName ||
				result.ClientEmail != tc.expected.ClientEmail ||
				result.PayoutID != tc.expected.PayoutID {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTransformDonationModelsToDTOs(t *testing.T) {
	testCases := map[string]struct {
		input    []*models.Donation
		expected []*dto.FormattedDonation
	}{
		"multipleDonations": {
			input: []*models.Donation{
				{
					ID:          "donation1",
					Created:     1700000000,
					Gross:       5000,
					Fee:         100,
					Net:         4900,
					ClientName:  "John Doe",
					ClientEmail: "john@example.com",
					PayoutID:    sql.NullString{String: "payout1", Valid: true},
				},
				{
					ID:          "donation2",
					Created:     1700000500,
					Gross:       10000,
					Fee:         500,
					Net:         9500,
					ClientName:  "Jane Doe",
					ClientEmail: "jane@example.com",
					PayoutID:    sql.NullString{Valid: false},
				},
			},
			expected: []*dto.FormattedDonation{
				{
					ID:          "donation1",
					Created:     "14 Nov 2023",
					Gross:       "50.00 lei",
					Fee:         "1.00 lei",
					Net:         "49.00 lei",
					ClientName:  "John Doe",
					ClientEmail: "john@example.com",
					PayoutID:    "payout1",
				},
				{
					ID:          "donation2",
					Created:     "14 Nov 2023",
					Gross:       "100.00 lei",
					Fee:         "5.00 lei",
					Net:         "95.00 lei",
					ClientName:  "Jane Doe",
					ClientEmail: "jane@example.com",
					PayoutID:    "",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformDonationModelsToDTOs(tc.input)

			if len(result) != len(tc.expected) {
				t.Fatalf("Expected %d donations, got %d", len(tc.expected), len(result))
			}

			for i, res := range result {
				if res.ID != tc.expected[i].ID ||
					res.Created != tc.expected[i].Created ||
					res.Gross != tc.expected[i].Gross ||
					res.Fee != tc.expected[i].Fee ||
					res.Net != tc.expected[i].Net ||
					res.ClientName != tc.expected[i].ClientName ||
					res.ClientEmail != tc.expected[i].ClientEmail ||
					res.PayoutID != tc.expected[i].PayoutID {
					t.Errorf("Expected %v, got %v", tc.expected[i], res)
				}
			}
		})
	}
}

func TestTransformDonationModelToPayoutReportItem(t *testing.T) {
	testCases := map[string]struct {
		input    *models.Donation
		expected *dto.PayoutReportItem
	}{
		"validDonation": {
			input: &models.Donation{
				ID:      "donation1",
				Created: 1700000000,
				Gross:   5000,
				Fee:     100,
				Net:     4900,
			},
			expected: &dto.PayoutReportItem{
				ID:          "donation1",
				Type:        "donation",
				Description: "",
				Created:     "14 Nov 2023",
				Gross:       "50.00 lei",
				Fee:         "1.00 lei",
				Net:         "49.00 lei",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformDonationModelToPayoutReportItem(tc.input)

			if result.ID != tc.expected.ID ||
				result.Type != tc.expected.Type ||
				result.Description != tc.expected.Description ||
				result.Created != tc.expected.Created ||
				result.Gross != tc.expected.Gross ||
				result.Fee != tc.expected.Fee ||
				result.Net != tc.expected.Net {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTransformDonationModelsToPayoutReportItems(t *testing.T) {
	testCases := map[string]struct {
		input    []*models.Donation
		expected []*dto.PayoutReportItem
	}{
		"multipleDonations": {
			input: []*models.Donation{
				{
					ID:      "donation1",
					Created: 1700000000,
					Gross:   5000,
					Fee:     100,
					Net:     4900,
				},
				{
					ID:      "donation2",
					Created: 1700000500,
					Gross:   10000,
					Fee:     500,
					Net:     9500,
				},
			},
			expected: []*dto.PayoutReportItem{
				{
					ID:          "donation1",
					Type:        "donation",
					Description: "",
					Created:     "14 Nov 2023",
					Gross:       "50.00 lei",
					Fee:         "1.00 lei",
					Net:         "49.00 lei",
				},
				{
					ID:          "donation2",
					Type:        "donation",
					Description: "",
					Created:     "14 Nov 2023",
					Gross:       "100.00 lei",
					Fee:         "5.00 lei",
					Net:         "95.00 lei",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformDonationModelsToPayoutReportItems(tc.input)

			if len(result) != len(tc.expected) {
				t.Fatalf("Expected %d donations, got %d", len(tc.expected), len(result))
			}

			for i, res := range result {
				if res.ID != tc.expected[i].ID ||
					res.Type != tc.expected[i].Type ||
					res.Description != tc.expected[i].Description ||
					res.Created != tc.expected[i].Created ||
					res.Gross != tc.expected[i].Gross ||
					res.Fee != tc.expected[i].Fee ||
					res.Net != tc.expected[i].Net {
					t.Errorf("Expected %v, got %v", tc.expected[i], res)
				}
			}
		})
	}
}

func TestTransformFeeModelToDTO(t *testing.T) {
	testCases := map[string]struct {
		input    *models.Fee
		expected *dto.FormattedFee
	}{
		"validFee": {
			input: &models.Fee{
				ID:          "fee1",
				Description: "Transaction fee",
				Created:     1700000000,
				Fee:         500,
			},
			expected: &dto.FormattedFee{
				ID:          "fee1",
				Description: "Transaction fee",
				Created:     "14 Nov 2023",
				Fee:         "5.00 lei",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformFeeModelToDTO(tc.input)

			if result.ID != tc.expected.ID ||
				result.Description != tc.expected.Description ||
				result.Created != tc.expected.Created ||
				result.Fee != tc.expected.Fee {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTransformFeeModelsToDTOs(t *testing.T) {
	testCases := map[string]struct {
		input    []*models.Fee
		expected []*dto.FormattedFee
	}{
		"multipleFees": {
			input: []*models.Fee{
				{ID: "fee1", Description: "Transaction fee", Created: 1700000000, Fee: 500},
				{ID: "fee2", Description: "Service fee", Created: 1700000500, Fee: 1000},
			},
			expected: []*dto.FormattedFee{
				{ID: "fee1", Description: "Transaction fee", Created: "14 Nov 2023", Fee: "5.00 lei"},
				{ID: "fee2", Description: "Service fee", Created: "14 Nov 2023", Fee: "10.00 lei"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformFeeModelsToDTOs(tc.input)

			if len(result) != len(tc.expected) {
				t.Fatalf("Expected %d fees, got %d", len(tc.expected), len(result))
			}

			for i, res := range result {
				if res.ID != tc.expected[i].ID ||
					res.Description != tc.expected[i].Description ||
					res.Created != tc.expected[i].Created ||
					res.Fee != tc.expected[i].Fee {
					t.Errorf("Expected %v, got %v", tc.expected[i], res)
				}
			}
		})
	}
}

func TestTransformFeeModelToPayoutReportItem(t *testing.T) {
	testCases := map[string]struct {
		input    *models.Fee
		expected *dto.PayoutReportItem
	}{
		"validFeeReportItem": {
			input: &models.Fee{
				ID:          "fee1",
				Description: "Transaction fee",
				Created:     1700000000,
				Fee:         500,
			},
			expected: &dto.PayoutReportItem{
				ID:          "fee1",
				Type:        "fee",
				Description: "Transaction fee",
				Created:     "14 Nov 2023",
				Gross:       "0.00 lei",
				Fee:         "5.00 lei",
				Net:         "5.00 lei",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformFeeModelToPayoutReportItem(tc.input)

			if result.ID != tc.expected.ID ||
				result.Type != tc.expected.Type ||
				result.Description != tc.expected.Description ||
				result.Created != tc.expected.Created ||
				result.Gross != tc.expected.Gross ||
				result.Fee != tc.expected.Fee ||
				result.Net != tc.expected.Net {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTransformFeeModelsToPayoutReportItems(t *testing.T) {
	testCases := map[string]struct {
		input    []*models.Fee
		expected []*dto.PayoutReportItem
	}{
		"multipleFeeReportItems": {
			input: []*models.Fee{
				{ID: "fee1", Description: "Transaction fee", Created: 1700000000, Fee: 500},
				{ID: "fee2", Description: "Service fee", Created: 1700000500, Fee: 1000},
			},
			expected: []*dto.PayoutReportItem{
				{ID: "fee1", Type: "fee", Description: "Transaction fee", Created: "14 Nov 2023", Gross: "0.00 lei", Fee: "5.00 lei", Net: "5.00 lei"},
				{ID: "fee2", Type: "fee", Description: "Service fee", Created: "14 Nov 2023", Gross: "0.00 lei", Fee: "10.00 lei", Net: "10.00 lei"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformFeeModelsToPayoutReportItems(tc.input)

			if len(result) != len(tc.expected) {
				t.Fatalf("Expected %d fees, got %d", len(tc.expected), len(result))
			}

			for i, res := range result {
				if res.ID != tc.expected[i].ID ||
					res.Type != tc.expected[i].Type ||
					res.Description != tc.expected[i].Description ||
					res.Created != tc.expected[i].Created ||
					res.Gross != tc.expected[i].Gross ||
					res.Fee != tc.expected[i].Fee ||
					res.Net != tc.expected[i].Net {
					t.Errorf("Expected %v, got %v", tc.expected[i], res)
				}
			}
		})
	}
}

func TestGetUnixTimestampsForMonth(t *testing.T) {
	testCases := map[string]struct {
		input    time.Time
		expected struct {
			monthStart int64
			monthEnd   int64
		}
	}{
		"validMonth": {
			input: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
			expected: struct {
				monthStart int64
				monthEnd   int64
			}{
				monthStart: 1725148800,
				monthEnd:   1727740799,
			},
		},
		"leapYearFebruary": {
			input: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected: struct {
				monthStart int64
				monthEnd   int64
			}{
				monthStart: 1706745600,
				monthEnd:   1709251199,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			monthStart, monthEnd := getUnixTimestampsForMonth(tc.input)

			if monthStart != tc.expected.monthStart {
				t.Errorf("Expected monthStart %d, got %d", tc.expected.monthStart, monthStart)
			}
			if monthEnd != tc.expected.monthEnd {
				t.Errorf("Expected monthEnd %d, got %d", tc.expected.monthEnd, monthEnd)
			}
		})
	}
}

func TestGetMonthDatesFromISO(t *testing.T) {
	testCases := map[string]struct {
		input    time.Time
		expected struct {
			monthStart   string
			monthEnd     string
			emissionDate string
		}
	}{
		"validMonth": {
			input: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
			expected: struct {
				monthStart   string
				monthEnd     string
				emissionDate string
			}{
				monthStart:   "1 Sep, 2024",
				monthEnd:     "30 Sep, 2024",
				emissionDate: "1 Oct, 2024",
			},
		},
		"leapYearFebruary": {
			input: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected: struct {
				monthStart   string
				monthEnd     string
				emissionDate string
			}{
				monthStart:   "1 Feb, 2024",
				monthEnd:     "29 Feb, 2024",
				emissionDate: "1 Mar, 2024",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			monthStart, monthEnd, emissionDate := getMonthDatesFromISO(tc.input)

			if monthStart != tc.expected.monthStart {
				t.Errorf("Expected monthStart %s, got %s", tc.expected.monthStart, monthStart)
			}
			if monthEnd != tc.expected.monthEnd {
				t.Errorf("Expected monthEnd %s, got %s", tc.expected.monthEnd, monthEnd)
			}
			if emissionDate != tc.expected.emissionDate {
				t.Errorf("Expected emissionDate %s, got %s", tc.expected.emissionDate, emissionDate)
			}
		})
	}
}

func TestValidateMonthString(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected time.Time
		err      bool
	}{
		"validDate": {
			input:    "2024-09",
			expected: time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
			err:      false,
		},
		"invalidDateFormat": {
			input:    "09-2024",
			expected: time.Time{},
			err:      true,
		},
		"nonExistentMonth": {
			input:    "2024-13",
			expected: time.Time{},
			err:      true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := validateMonthString(tc.input)

			if tc.err {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
				if !result.Equal(tc.expected) {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

func TestTransformToMonthlyReportData(t *testing.T) {
	payoutModels := []*models.Payout{
		{ID: "payout1", Created: 1700000000, Gross: 10000, Fee: 1000, Net: 9000},
		{ID: "payout2", Created: 1700000001, Gross: 20000, Fee: 2000, Net: 18000},
	}

	testCases := map[string]struct {
		date   time.Time
		gross  uint32
		fee    uint32
		net    uint32
		expect *dto.MonthlyReportData
	}{
		"validInput": {
			date:  time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC),
			gross: 30000,
			fee:   3000,
			net:   27000,
			expect: dto.NewMonthlyReportData(
				"1 Sep, 2024",
				"30 Sep, 2024",
				"1 Oct, 2024",
				"300.00 lei",
				"30.00 lei",
				"270.00 lei",
				transformPayoutModelsToDTOs(payoutModels),
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformToMonthlyReportData(tc.date, tc.gross, tc.fee, tc.net, payoutModels)

			if result.MonthStart != tc.expect.MonthStart {
				t.Errorf("Expected MonthStart %s, got %s", tc.expect.MonthStart, result.MonthStart)
			}
			if result.MonthEnd != tc.expect.MonthEnd {
				t.Errorf("Expected MonthEnd %s, got %s", tc.expect.MonthEnd, result.MonthEnd)
			}
			if result.EmissionDate != tc.expect.EmissionDate {
				t.Errorf("Expected EmissionDate %s, got %s", tc.expect.EmissionDate, result.EmissionDate)
			}
			if result.Gross != tc.expect.Gross {
				t.Errorf("Expected Gross %s, got %s", tc.expect.Gross, result.Gross)
			}
			if result.Fee != tc.expect.Fee {
				t.Errorf("Expected Fee %s, got %s", tc.expect.Fee, result.Fee)
			}
			if result.Net != tc.expect.Net {
				t.Errorf("Expected Net %s, got %s", tc.expect.Net, result.Net)
			}
			if len(result.Payouts) != len(tc.expect.Payouts) {
				t.Errorf("Expected %d payouts, got %d", len(tc.expect.Payouts), len(result.Payouts))
			}
		})
	}
}

func TestMonthlyReportSum(t *testing.T) {
	testCases := map[string]struct {
		payouts              []*models.Payout
		expectedGross        uint32
		expectedFee          uint32
		expectedNet          uint32
		expectError          bool
		expectedErrorMessage string
	}{
		"validPayouts": {
			payouts: []*models.Payout{
				{Gross: 10000, Fee: 1000, Net: 9000},
				{Gross: 20000, Fee: 2000, Net: 18000},
			},
			expectedGross: 30000,
			expectedFee:   3000,
			expectedNet:   27000,
			expectError:   false,
		},
		"mismatchNet": {
			payouts: []*models.Payout{
				{Gross: 10000, Fee: 1000, Net: 8000},
				{Gross: 20000, Fee: 2000, Net: 18000},
			},
			expectedGross:        0,
			expectedFee:          0,
			expectedNet:          0,
			expectError:          true,
			expectedErrorMessage: "monthly gross-fee 27000 != net 26000",
		},
		"noPayouts": {
			payouts:       nil,
			expectedGross: 0,
			expectedFee:   0,
			expectedNet:   0,
			expectError:   false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gross, fee, net, err := monthlyReportSum(tc.payouts)

			if gross != tc.expectedGross {
				t.Errorf("Expected gross %d, got %d", tc.expectedGross, gross)
			}
			if fee != tc.expectedFee {
				t.Errorf("Expected fee %d, got %d", tc.expectedFee, fee)
			}
			if net != tc.expectedNet {
				t.Errorf("Expected net %d, got %d", tc.expectedNet, net)
			}
			if tc.expectError && err == nil {
				t.Error("Expected an error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectError && err != nil {
				if err.Error() != tc.expectedErrorMessage {
					t.Errorf("Expected error message %q, got %q", tc.expectedErrorMessage, err.Error())
				}
			}
		})
	}
}
