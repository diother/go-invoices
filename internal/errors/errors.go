package errors

// Charge-related errors
const (
	ErrChargeMissing             = "charge object is nil"
	ErrChargeStatusInvalid       = "charge status is not succeeded"
	ErrChargeBillingMissing      = "charge billing details object is nil"
	ErrChargeBillingNameMissing  = "charge billing details name is missing"
	ErrChargeBillingEmailMissing = "charge billing details email is missing"
)

// Payout-related errors
const (
	ErrPayoutMissing       = "payout object is nil"
	ErrPayoutIDMissing     = "payout ID is missing"
	ErrPayoutStatusInvalid = "payout status is not paid"
)

// Payout list validation errors
// NOTE: in work
const (
	ErrPayoutListTransactionMismatch = "payout amount does not match total charges minus fees"
	ErrPayoutListTransactionInvalid  = "transaction list expected at least 2 transactions"
)

// General transaction-related errors
const (
	ErrTransactionMissing        = "transaction object is nil"
	ErrTransactionIDMissing      = "transaction ID is missing"
	ErrTransactionCreatedInvalid = "transaction creation date is invalid"
)

// Charge transaction-related errors
const (
	ErrChargeTransactionTypeInvalid     = "transaction is not of type charge"
	ErrChargeTransactionAmountInvalid   = "charge transaction amount is missing, zero, or negative"
	ErrChargeTransactionFeeInvalid      = "charge transaction fee is missing, zero, or negative"
	ErrChargeTransactionNetInvalid      = "charge transaction net is missing, zero, or negative"
	ErrChargeTransactionSourceMissing   = "charge transaction source is missing"
	ErrChargeTransactionSourceIDMissing = "charge transaction source ID is missing"
)

// Payout transaction-related errors
const (
	ErrPayoutTransactionTypeInvalid   = "transaction is not of type payout"
	ErrPayoutTransactionAmountInvalid = "payout transaction amount is missing, zero, or positive"
	ErrPayoutTransactionFeeInvalid    = "payout transaction fee is not zero"
	ErrPayoutTransactionNetInvalid    = "payout transaction net is missing, zero, or positive"
)

// Fee transaction-related errors
const (
	ErrFeeTransactionTypeInvalid   = "transaction is not of type stripe_fee"
	ErrFeeTransactionAmountInvalid = "fee transaction amount is missing, zero, or positive"
	ErrFeeTransactionFeeInvalid    = "fee transaction fee is not zero"
	ErrFeeTransactionNetInvalid    = "fee transaction net is missing, zero, or positive"
)
