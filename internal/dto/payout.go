package dto

type FormattedPayout struct {
	ID        string
	Created   string
	Gross     string
	Fee       string
	Net       string
	Donations []*FormattedDonation
	Fees      []*FormattedFee
}

func NewFormattedPayout(id, created, gross, fee, net string, donations []*FormattedDonation, fees []*FormattedFee) *FormattedPayout {
	return &FormattedPayout{
		ID:        id,
		Created:   created,
		Gross:     gross,
		Fee:       fee,
		Net:       net,
		Donations: donations,
		Fees:      fees,
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
	ID          string
	Type        string
	Description string
	Created     string
	Gross       string
	Fee         string
	Net         string
}

func NewPayoutReportItem(id, itemType, description, created, gross, fee, net string) *PayoutReportItem {
	return &PayoutReportItem{
		ID:          id,
		Type:        itemType,
		Description: description,
		Created:     created,
		Gross:       gross,
		Fee:         fee,
		Net:         net,
	}
}

type MonthlyReportData struct {
	MonthStart   string
	MonthEnd     string
	EmissionDate string
	Gross        string
	Fee          string
	Net          string
	Payouts      []*FormattedPayout
}

func NewMonthlyReportData(monthStart, monthEnd, emissionDate, gross, fee, net string, payouts []*FormattedPayout) *MonthlyReportData {
	return &MonthlyReportData{
		MonthStart:   monthStart,
		MonthEnd:     monthEnd,
		EmissionDate: emissionDate,
		Gross:        gross,
		Fee:          fee,
		Net:          net,
		Payouts:      payouts,
	}
}

type MonthlyReportView struct {
	Date    string
	Gross   string
	Fee     string
	Net     string
	Payouts []*FormattedPayout
}

func NewMonthlyReportView(date, gross, fee, net string, payouts []*FormattedPayout) *MonthlyReportView {
	return &MonthlyReportView{
		Date:    date,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
		Payouts: payouts,
	}
}
