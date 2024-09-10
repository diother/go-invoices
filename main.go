package main

import (
	"log"
	"path/filepath"

	"github.com/signintech/gopdf"
)

func generatePDF(filename string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Roboto", "./static/Roboto-Regular.ttf")
	if err != nil {
		return err
	}
	pdf.SetFont("Roboto", "", 14)

	pdf.Cell(nil, "Salutare, Dan")

	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, filename+".pdf")
	return pdf.WritePdf(pdfFile)
}

func main() {
	err := generatePDF("test")
	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
