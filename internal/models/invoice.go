package models

type Invoice struct {
	TransactionId string
	IssueDate     string
	ClientName    string
	ProductName   string
	UnitPrice     string
	Total         string
}
