package main

import (
	"testing"
)

func BenchmarkGeneratePDF(b *testing.B) {
	inv := Invoice{
		ClientName:    "Ungureanu Daniel",
		IssueDate:     "12 Aug, 2024",
		TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
		ProductName:   "Donație unică de 10 lei",
		UnitPrice:     10.00,
		Total:         10.00,
	}

	for i := 0; i < b.N; i++ {
		err := generateInvoicePdf(&inv)
		if err != nil {
			b.Fatalf("Error generating PDF: %v", err)
		}
	}
}
