package models

type Invoice struct {
	TransactionId string
	IssueDate     string
	ClientName    string
	ProductName   string
	UnitPrice     float64
	Total         float64
}
