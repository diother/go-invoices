package services

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/diother/go-invoices/internal/constants"
	"github.com/diother/go-invoices/internal/models"
	"github.com/stripe/stripe-go/v79"
)

func TestValidatePayout(t *testing.T) {
	testCases := map[string]struct {
		input    *stripe.Payout
		expected string
	}{
		"validPayout":   {&stripe.Payout{ID: "po_123456789", Status: "paid"}, ""},
		"payoutMissing": {nil, constants.ErrPayoutMissing},
		"statusInvalid": {&stripe.Payout{ID: "po_123456789", Status: "pending"}, constants.ErrPayoutStatusInvalid},
		"IDMissing":     {&stripe.Payout{Status: "paid"}, constants.ErrPayoutIDMissing},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validatePayout(tc.input)
			if tc.expected == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expected != "" && (err == nil || err.Error() != tc.expected) {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestValidateCharge(t *testing.T) {
	testCases := map[string]struct {
		input    *stripe.Charge
		expected string
	}{
		"validCharge": {
			&stripe.Charge{
				ID:     "ch_123456789",
				Status: "succeeded",
				BillingDetails: &stripe.ChargeBillingDetails{
					Name:  "John Doe",
					Email: "john.doe@example.com",
				},
				BalanceTransaction: &stripe.BalanceTransaction{ID: "txn_123456"},
			},
			"",
		},
		"chargeMissing":  {nil, constants.ErrChargeMissing},
		"statusInvalid":  {&stripe.Charge{ID: "ch_123456789", Status: "failed"}, constants.ErrChargeStatusInvalid},
		"billingMissing": {&stripe.Charge{ID: "ch_123456789", Status: "succeeded"}, constants.ErrChargeBillingMissing},
		"billingNameMissing": {&stripe.Charge{
			ID:     "ch_123456789",
			Status: "succeeded",
			BillingDetails: &stripe.ChargeBillingDetails{
				Email: "john.doe@example.com",
			},
		}, constants.ErrChargeBillingNameMissing},
		"billingEmailMissing": {&stripe.Charge{
			ID:     "ch_123456789",
			Status: "succeeded",
			BillingDetails: &stripe.ChargeBillingDetails{
				Name: "John Doe",
			},
		}, constants.ErrChargeBillingEmailMissing},
		"transactionMissing": {&stripe.Charge{
			ID:     "ch_123456789",
			Status: "succeeded",
			BillingDetails: &stripe.ChargeBillingDetails{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		}, constants.ErrTransactionMissing},
		"transactionIDMissing": {&stripe.Charge{
			ID:     "ch_123456789",
			Status: "succeeded",
			BillingDetails: &stripe.ChargeBillingDetails{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			BalanceTransaction: &stripe.BalanceTransaction{},
		}, constants.ErrTransactionIDMissing},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateCharge(tc.input)
			if tc.expected == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expected != "" && (err == nil || err.Error() != tc.expected) {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestValidatePayoutTransaction(t *testing.T) {
	testCases := map[string]struct {
		input    *stripe.BalanceTransaction
		expected string
	}{
		"validTransaction": {
			&stripe.BalanceTransaction{
				ID:      "txn_123456",
				Type:    "payout",
				Created: 1234567890,
				Amount:  -1000,
				Fee:     0,
				Net:     -1000,
			},
			"",
		},
		"transactionMissing": {nil, constants.ErrTransactionMissing},
		"typeInvalid":        {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge"}, constants.ErrPayoutTransactionTypeInvalid},
		"IDMissing":          {&stripe.BalanceTransaction{Type: "payout"}, constants.ErrTransactionIDMissing},
		"createdInvalid":     {&stripe.BalanceTransaction{ID: "txn_123456", Type: "payout", Created: 0}, constants.ErrTransactionCreatedInvalid},
		"amountInvalid":      {&stripe.BalanceTransaction{ID: "txn_123456", Type: "payout", Created: 1234567890, Amount: 0}, constants.ErrPayoutTransactionAmountInvalid},
		"feeInvalid":         {&stripe.BalanceTransaction{ID: "txn_123456", Type: "payout", Created: 1234567890, Amount: -1000, Fee: 1}, constants.ErrPayoutTransactionFeeInvalid},
		"netInvalid":         {&stripe.BalanceTransaction{ID: "txn_123456", Type: "payout", Created: 1234567890, Amount: -1000, Net: 0}, constants.ErrPayoutTransactionNetInvalid},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validatePayoutTransaction(tc.input)
			if tc.expected == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expected != "" && (err == nil || err.Error() != tc.expected) {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestValidateChargeTransaction(t *testing.T) {
	testCases := map[string]struct {
		input    *stripe.BalanceTransaction
		expected string
	}{
		"validTransaction": {
			&stripe.BalanceTransaction{
				ID:      "txn_123456",
				Type:    "charge",
				Created: 1234567890,
				Amount:  1000,
				Fee:     100,
				Net:     900,
				Source:  &stripe.BalanceTransactionSource{ID: "src_123456"},
			},
			"",
		},
		"transactionMissing": {nil, constants.ErrTransactionMissing},
		"typeInvalid":        {&stripe.BalanceTransaction{ID: "txn_123456", Type: "payout"}, constants.ErrChargeTransactionTypeInvalid},
		"IDMissing":          {&stripe.BalanceTransaction{Type: "charge"}, constants.ErrTransactionIDMissing},
		"createdInvalid":     {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 0}, constants.ErrTransactionCreatedInvalid},
		"amountInvalid":      {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 0}, constants.ErrChargeTransactionAmountInvalid},
		"feeInvalid":         {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 1000, Fee: 0}, constants.ErrChargeTransactionFeeInvalid},
		"netInvalid":         {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 1000, Fee: 100, Net: 0}, constants.ErrChargeTransactionNetInvalid},
		"sourceMissing":      {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 1000, Fee: 100, Net: 900}, constants.ErrChargeTransactionSourceMissing},
		"sourceIDMissing":    {&stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 1000, Fee: 100, Net: 900, Source: &stripe.BalanceTransactionSource{}}, constants.ErrChargeTransactionSourceIDMissing},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateChargeTransaction(tc.input)
			if tc.expected == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expected != "" && (err == nil || err.Error() != tc.expected) {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestValidateFeeTransaction(t *testing.T) {
	testCases := map[string]struct {
		input    *stripe.BalanceTransaction
		expected string
	}{
		"validTransaction": {
			&stripe.BalanceTransaction{
				ID:          "txn_fee_123456",
				Type:        "stripe_fee",
				Description: "Billing",
				Created:     1234567890,
				Amount:      -100,
				Fee:         0,
				Net:         -100,
			},
			"",
		},
		"transactionMissing": {nil, constants.ErrTransactionMissing},
		"typeInvalid":        {&stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "charge"}, constants.ErrFeeTransactionTypeInvalid},
		"IDMissing":          {&stripe.BalanceTransaction{Type: "stripe_fee"}, constants.ErrTransactionIDMissing},
		"descriptionMissing": {&stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee"}, constants.ErrFeeTransactionDescriptionMissing},
		"createdInvalid":     {&stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee", Description: "Billing"}, constants.ErrTransactionCreatedInvalid},
		"amountInvalid":      {&stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee", Description: "Billing", Created: 1234567890, Amount: 0}, constants.ErrFeeTransactionAmountInvalid},
		"feeInvalid":         {&stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee", Description: "Billing", Created: 1234567890, Amount: -100, Fee: 1}, constants.ErrFeeTransactionFeeInvalid},
		"netInvalid":         {&stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee", Description: "Billing", Created: 1234567890, Amount: -100, Fee: 0, Net: 1}, constants.ErrFeeTransactionNetInvalid},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateFeeTransaction(tc.input)
			if tc.expected == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expected != "" && (err == nil || err.Error() != tc.expected) {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestValidateRelatedTransactions(t *testing.T) {
	validPayout := &stripe.BalanceTransaction{ID: "txn_123456", Type: "payout", Created: 1234567890, Amount: -1000, Fee: 0, Net: -1000}
	validCharge := &stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 1000, Fee: 100, Net: 900, Source: &stripe.BalanceTransactionSource{ID: "src_123456"}}
	validFee := &stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee", Description: "Billing", Created: 1234567890, Amount: -100, Fee: 0, Net: -100}

	testCases := map[string]struct {
		input    []*stripe.BalanceTransaction
		expected string
	}{
		"validTransactions": {
			input:    []*stripe.BalanceTransaction{validPayout, validCharge, validFee},
			expected: "",
		},
		"insufficientTransactions": {
			input:    []*stripe.BalanceTransaction{validPayout},
			expected: constants.ErrPayoutListInsufficientTransactions,
		},
		"payoutTransactionInvalid": {
			input:    []*stripe.BalanceTransaction{validCharge, validCharge},
			expected: constants.ErrPayoutListPayoutTransactionInvalid,
		},
		"relatedTransactionInvalid": {
			input:    []*stripe.BalanceTransaction{validPayout, {Type: "charge", ID: ""}},
			expected: constants.ErrPayoutListRelatedTransactionInvalid,
		},
		"unexpectedType": {
			input:    []*stripe.BalanceTransaction{validPayout, {Type: "unexpected"}},
			expected: constants.ErrPayoutListUnexpectedTransaction,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateRelatedTransactions(tc.input)
			if tc.expected == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expected != "" && (err == nil || !strings.Contains(err.Error(), tc.expected)) {
				t.Errorf("Expected error containing: %v, got: %v", tc.expected, err)
			}
		})
	}
}

func TestValidateMatchingSums(t *testing.T) {
	validPayout := &stripe.BalanceTransaction{ID: "txn_123456", Type: "payout", Created: 1234567890, Amount: -800, Fee: 0, Net: -800}
	validCharge := &stripe.BalanceTransaction{ID: "txn_123456", Type: "charge", Created: 1234567890, Amount: 1000, Fee: 100, Net: 900, Source: &stripe.BalanceTransactionSource{ID: "src_123456"}}
	validFee := &stripe.BalanceTransaction{ID: "txn_fee_123456", Type: "stripe_fee", Description: "Billing", Created: 1234567890, Amount: -100, Fee: 0, Net: -100}

	testCases := map[string]struct {
		input         []*stripe.BalanceTransaction
		expectedErr   string
		expectedGross int64
		expectedFee   int64
		expectedNet   int64
	}{
		"validRelatedTransactions": {
			input:         []*stripe.BalanceTransaction{validPayout, validCharge, validFee},
			expectedErr:   "",
			expectedGross: 1000,
			expectedFee:   200,
			expectedNet:   800,
		},
		"payoutMismatch": {
			input:       []*stripe.BalanceTransaction{validPayout, validCharge, validFee, validFee},
			expectedErr: constants.ErrPayoutListSumMismatch,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gross, fee, net, err := validateMatchingSums(tc.input)
			if tc.expectedErr == "" && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedErr != "" && (err == nil || !strings.Contains(err.Error(), tc.expectedErr)) {
				t.Errorf("Expected error containing: %v, got: %v", tc.expectedErr, err)
			}
			if err == nil {
				if gross != tc.expectedGross {
					t.Errorf("Expected gross: %d, got: %d", tc.expectedGross, gross)
				}
				if fee != tc.expectedFee {
					t.Errorf("Expected fee: %d, got: %d", tc.expectedFee, fee)
				}
				if net != tc.expectedNet {
					t.Errorf("Expected net: %d, got: %d", tc.expectedNet, net)
				}
			}
		})
	}
}

func TestTransformNoPayoutDonationDTOToModel(t *testing.T) {
	testCases := map[string]struct {
		transaction *stripe.BalanceTransaction
		charge      *stripe.Charge
		expected    *models.Donation
	}{
		"validData": {
			transaction: &stripe.BalanceTransaction{
				ID:      "txn_123456",
				Created: 1234567890,
				Amount:  2500,
				Fee:     500,
				Net:     2000,
			},
			charge: &stripe.Charge{
				BillingDetails: &stripe.ChargeBillingDetails{
					Name:  "John Doe",
					Email: "john.doe@example.com",
				},
			},
			expected: models.NewDonation(
				"txn_123456",
				uint64(1234567890),
				uint32(2500),
				uint32(500),
				uint32(2000),
				"John Doe",
				"john.doe@example.com",
				sql.NullString{Valid: false},
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformNoPayoutDonationDTOToModel(tc.transaction, tc.charge)

			if result.ID != tc.expected.ID {
				t.Errorf("Expected ID %v, got %v", tc.expected.ID, result.ID)
			}
			if result.Created != tc.expected.Created {
				t.Errorf("Expected Created %v, got %v", tc.expected.Created, result.Created)
			}
			if result.Gross != tc.expected.Gross {
				t.Errorf("Expected Gross %v, got %v", tc.expected.Gross, result.Gross)
			}
			if result.Fee != tc.expected.Fee {
				t.Errorf("Expected Fee %v, got %v", tc.expected.Fee, result.Fee)
			}
			if result.Net != tc.expected.Net {
				t.Errorf("Expected Net %v, got %v", tc.expected.Net, result.Net)
			}
			if result.ClientName != tc.expected.ClientName {
				t.Errorf("Expected ClientName %v, got %v", tc.expected.ClientName, result.ClientName)
			}
			if result.ClientEmail != tc.expected.ClientEmail {
				t.Errorf("Expected ClientEmail %v, got %v", tc.expected.ClientEmail, result.ClientEmail)
			}
			if result.PayoutID != tc.expected.PayoutID {
				t.Errorf("Expected PayoutID %v, got %v", tc.expected.PayoutID, result.PayoutID)
			}
		})
	}
}

func TestTransformUpdateDonationDTOToModel(t *testing.T) {
	testCases := map[string]struct {
		transactionID string
		payoutID      string
		expected      *models.Donation
	}{
		"validData": {
			transactionID: "txn_7894561",
			payoutID:      "txn_789456",
			expected: models.NewDonation(
				"txn_7894561",
				0, 0, 0, 0, "", "",
				sql.NullString{String: "txn_789456", Valid: true},
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformUpdateDonationDTOToModel(tc.transactionID, tc.payoutID)

			if result.ID != tc.expected.ID {
				t.Errorf("Expected ID %v, got %v", tc.expected.ID, result.ID)
			}
			if result.PayoutID != tc.expected.PayoutID {
				t.Errorf("Expected PayoutID %v, got %v", tc.expected.PayoutID, result.PayoutID)
			}
			if result.Created != tc.expected.Created {
				t.Errorf("Expected Created %v, got %v", tc.expected.Created, result.Created)
			}
			if result.Gross != tc.expected.Gross {
				t.Errorf("Expected Gross %v, got %v", tc.expected.Gross, result.Gross)
			}
			if result.Fee != tc.expected.Fee {
				t.Errorf("Expected Fee %v, got %v", tc.expected.Fee, result.Fee)
			}
			if result.Net != tc.expected.Net {
				t.Errorf("Expected Net %v, got %v", tc.expected.Net, result.Net)
			}
			if result.ClientName != tc.expected.ClientName {
				t.Errorf("Expected ClientName %v, got %v", tc.expected.ClientName, result.ClientName)
			}
			if result.ClientEmail != tc.expected.ClientEmail {
				t.Errorf("Expected ClientEmail %v, got %v", tc.expected.ClientEmail, result.ClientEmail)
			}
		})
	}
}

func TestTransformDonationDTOToModel(t *testing.T) {
	testCases := map[string]struct {
		transaction *stripe.BalanceTransaction
		charge      *stripe.Charge
		payoutID    string
		expected    *models.Donation
	}{
		"validData": {
			transaction: &stripe.BalanceTransaction{
				ID:      "txn_987654",
				Created: 1627849100,
				Amount:  5000,
				Fee:     1000,
				Net:     4000,
			},
			charge: &stripe.Charge{
				BillingDetails: &stripe.ChargeBillingDetails{
					Name:  "Jane Doe",
					Email: "jane.doe@example.com",
				},
			},
			payoutID: "po_321654",
			expected: models.NewDonation(
				"txn_987654",
				uint64(1627849100),
				uint32(5000),
				uint32(1000),
				uint32(4000),
				"Jane Doe",
				"jane.doe@example.com",
				sql.NullString{String: "po_321654", Valid: true},
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformDonationDTOToModel(tc.transaction, tc.charge, tc.payoutID)

			if result.ID != tc.expected.ID {
				t.Errorf("Expected ID %v, got %v", tc.expected.ID, result.ID)
			}
			if result.Created != tc.expected.Created {
				t.Errorf("Expected Created %v, got %v", tc.expected.Created, result.Created)
			}
			if result.Gross != tc.expected.Gross {
				t.Errorf("Expected Gross %v, got %v", tc.expected.Gross, result.Gross)
			}
			if result.Fee != tc.expected.Fee {
				t.Errorf("Expected Fee %v, got %v", tc.expected.Fee, result.Fee)
			}
			if result.Net != tc.expected.Net {
				t.Errorf("Expected Net %v, got %v", tc.expected.Net, result.Net)
			}
			if result.ClientName != tc.expected.ClientName {
				t.Errorf("Expected ClientName %v, got %v", tc.expected.ClientName, result.ClientName)
			}
			if result.ClientEmail != tc.expected.ClientEmail {
				t.Errorf("Expected ClientEmail %v, got %v", tc.expected.ClientEmail, result.ClientEmail)
			}
			if result.PayoutID != tc.expected.PayoutID {
				t.Errorf("Expected PayoutID %v, got %v", tc.expected.PayoutID, result.PayoutID)
			}
		})
	}
}

func TestTransformPayoutDTOToModel(t *testing.T) {
	testCases := map[string]struct {
		transaction *stripe.BalanceTransaction
		gross       int64
		fee         int64
		net         int64
		expected    *models.Payout
	}{
		"validData": {
			transaction: &stripe.BalanceTransaction{
				ID:      "txn_payout_123456",
				Created: 1627849100,
			},
			gross: 10000,
			fee:   500,
			net:   9500,
			expected: models.NewPayout(
				"txn_payout_123456",
				uint64(1627849100),
				uint32(10000),
				uint32(500),
				uint32(9500),
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformPayoutDTOToModel(tc.transaction, tc.gross, tc.fee, tc.net)

			if result.ID != tc.expected.ID {
				t.Errorf("Expected ID %v, got %v", tc.expected.ID, result.ID)
			}
			if result.Created != tc.expected.Created {
				t.Errorf("Expected Created %v, got %v", tc.expected.Created, result.Created)
			}
			if result.Gross != tc.expected.Gross {
				t.Errorf("Expected Gross %v, got %v", tc.expected.Gross, result.Gross)
			}
			if result.Fee != tc.expected.Fee {
				t.Errorf("Expected Fee %v, got %v", tc.expected.Fee, result.Fee)
			}
			if result.Net != tc.expected.Net {
				t.Errorf("Expected Net %v, got %v", tc.expected.Net, result.Net)
			}
		})
	}
}

func TestTransformFeeDTOToModel(t *testing.T) {
	testCases := map[string]struct {
		transaction *stripe.BalanceTransaction
		payoutID    string
		expected    *models.Fee
	}{
		"validData": {
			transaction: &stripe.BalanceTransaction{
				ID:          "txn_fee_789456",
				Description: "billing",
				Created:     1627849100,
				Amount:      -1000,
			},
			payoutID: "po_321654",
			expected: models.NewFee(
				"txn_fee_789456",
				"billing",
				uint64(1627849100),
				uint32(1000),
				sql.NullString{String: "po_321654", Valid: true},
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := transformFeeDTOToModel(tc.transaction, tc.payoutID)

			if result.ID != tc.expected.ID {
				t.Errorf("Expected ID %v, got %v", tc.expected.ID, result.ID)
			}
			if result.Created != tc.expected.Created {
				t.Errorf("Expected Created %v, got %v", tc.expected.Created, result.Created)
			}
			if result.Fee != tc.expected.Fee {
				t.Errorf("Expected Fee %v, got %v", tc.expected.Fee, result.Fee)
			}
			if result.PayoutID != tc.expected.PayoutID {
				t.Errorf("Expected PayoutID %v, got %v", tc.expected.PayoutID, result.PayoutID)
			}
		})
	}
}
