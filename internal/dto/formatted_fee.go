package dto

type FormattedFee struct {
	ID      string
	Created string
	Fee     string
}

func NewFormattedFee(id, created, fee string) *FormattedFee {
	return &FormattedFee{
		ID:      id,
		Created: created,
		Fee:     fee,
	}
}
