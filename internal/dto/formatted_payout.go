package dto

type FormattedPayout struct {
	ID      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func NewFormattedPayout(id, created, gross, fee, net string) *FormattedPayout {
	return &FormattedPayout{
		ID:      id,
		Created: created,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}
