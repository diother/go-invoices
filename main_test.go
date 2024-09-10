package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePDF(t *testing.T) {
	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, "test_hello_world.pdf")

	defer os.Remove(pdfFile)

	err := generatePDF("test_hello_world")
	if err != nil {
		t.Fatalf("Error generating PDF: %v", err)
	}

	if _, err := os.Stat(pdfFile); os.IsNotExist(err) {
		t.Fatalf("PDF file does not exist: %v", err)
	}
}
