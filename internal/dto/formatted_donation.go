package dto

type FormattedDonation struct {
	ID          string
	Created     string
	Gross       string
	Fee         string
	Net         string
	ClientName  string
	ClientEmail string
	PayoutID    string
}

func NewFormattedDonation(id, created, gross, fee, net, clientName, clientEmail, payoutID string) *FormattedDonation {
	return &FormattedDonation{
		ID:          id,
		Created:     created,
		Gross:       gross,
		Fee:         fee,
		Net:         net,
		ClientName:  clientName,
		ClientEmail: clientEmail,
		PayoutID:    payoutID,
	}
}
