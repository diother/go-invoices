package models

type Payout struct {
	IssueDate  string
	PayoutDate string
	PayoutID   string
	Gross      string
	StripeFees string
	Total      string
	Items      []PayoutTransaction
}

type MonthlyPayout struct {
	IssueDate    string
	ReportPeriod string
	Gross        string
	StripeFees   string
	Total        string
	Items        []Payout
}

type PayoutTransaction struct {
	ProductName   string
	TransactionId string
	Gross         string
	StripeFee     string
	Total         string
}
