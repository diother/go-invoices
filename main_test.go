package main

import (
	"testing"

	"github.com/diother/go-invoices/models"
	"github.com/diother/go-invoices/views"
)

func BenchmarkGeneratePDF(b *testing.B) {
	inv := models.Invoice{
		ClientName:    "Ungureanu Daniel",
		IssueDate:     "12 Aug, 2024",
		TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
		ProductName:   "Donație unică de 10 lei",
		UnitPrice:     10.00,
		Total:         10.00,
	}
	adapter := views.InvoiceView{Invoice: &inv}

	for i := 0; i < b.N; i++ {
		err := adapter.GenerateDocument(&inv)
		if err != nil {
			b.Fatalf("Error generating PDF: %v", err)
		}
	}
}
