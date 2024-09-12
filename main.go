package main

import (
	"log"

	"github.com/diother/go-invoices/models"
	"github.com/diother/go-invoices/views"
)

func main() {
	inv := models.Invoice{
		ClientName:    "Ungureanu Daniel",
		IssueDate:     "12 Aug, 2024",
		TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
		ProductName:   "Donație unică de 10 lei",
		UnitPrice:     10,
		Total:         10,
	}
	err := views.GenerateInvoicePdf(&inv)
	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
