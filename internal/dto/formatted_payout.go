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

type PayoutReportData struct {
	Payout *FormattedPayout
	Items  []*PayoutReportItem
}

func NewPayoutReportData(payout *FormattedPayout, items []*PayoutReportItem) *PayoutReportData {
	return &PayoutReportData{
		Payout: payout,
		Items:  items,
	}
}

type PayoutReportItem struct {
	ID      string
	Type    string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func NewPayoutReportItem(id, itemType, created, gross, fee, net string) *PayoutReportItem {
	return &PayoutReportItem{
		ID:      id,
		Type:    itemType,
		Created: created,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}
