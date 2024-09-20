package errors

const (
	ErrChargeStatusFailed          = "Charge status is not 'succeeded'"
	ErrClientNameMissing           = "Client name is missing"
	ErrClientEmailMissing          = "Client email is missing"
	ErrBalanceTransactionIDMissing = "Balance transaction ID is missing"
	ErrTransactionIDMissing        = "Transaction ID is missing"
	ErrTransactionCreatedMissing   = "Transaction creation date is missing"
	ErrTransactionAmountMissing    = "Transaction amount is missing or zero"
	ErrTransactionFeeMissing       = "Transaction fee is missing or zero"
	ErrTransactionNetMissing       = "Transaction net amount is missing or zero"
	ErrTransactionPayoutFailed     = "Transaction is not of type payout"
	ErrTransactionChargeFailed     = "Transaction is not of type charge"
	ErrTransactionListMissing      = "Expected at least 2 transactions"
	ErrTransactionPayoutMismatch   = "Payout amount does not match total charges minus fees"
	ErrTransactionChargeIDMissing  = "Charge transaction source ID is missing"
)
