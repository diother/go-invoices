package main

import (
	"log"

	"github.com/diother/go-invoices/models"
	"github.com/diother/go-invoices/views"
)

func main() {

	pay := models.Payout{
		ClientName:    "Ungureanu Daniel",
		IssueDate:     "12 Aug, 2024",
		TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
		ProductName:   "Donație unică de 10 lei",
		UnitPrice:     10,
		Total:         10,
	}
	payout := views.PayoutView{Payout: &pay}
	err := payout.GenerateDocument(&pay)

	// inv := models.Invoice{
	// 	ClientName:    "Ungureanu Daniel",
	// 	IssueDate:     "12 Aug, 2024",
	// 	TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
	// 	ProductName:   "Donație unică de 10 lei",
	// 	UnitPrice:     10,
	// 	Total:         10,
	// }
	// invoice := views.InvoiceView{Invoice: &inv}
	// err := invoice.GenerateDocument(&inv)

	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
