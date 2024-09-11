package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePDF(t *testing.T) {
	inv := Invoice{
		ClientName:    "Ungureanu Daniel",
		IssueDate:     "12 Aug, 2024",
		TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
		ProductName:   "Donație unică de 10 lei",
		UnitPrice:     10.00,
		Total:         10.00,
	}
	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, "test_hello_world.pdf")

	defer os.Remove(pdfFile)

	err := generatePDF(inv)
	if err != nil {
		t.Fatalf("Error generating PDF: %v", err)
	}

	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		t.Fatalf("PDF file does not exist: %v", err)
	}
}
