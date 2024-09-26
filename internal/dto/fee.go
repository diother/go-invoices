package dto

type FormattedFee struct {
	ID          string
	Description string
	Created     string
	Fee         string
}

func NewFormattedFee(id, description, created, fee string) *FormattedFee {
	return &FormattedFee{
		ID:          id,
		Description: description,
		Created:     created,
		Fee:         fee,
	}
}
