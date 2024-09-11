package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/signintech/gopdf"
)

type Invoice struct {
	TransactionId string
	IssueDate     string
	ClientName    string
	ProductName   string
	UnitPrice     float64
	Total         float64
}

func rightAlign(pdf *gopdf.GoPdf, xEnd, y float64, text string) {
	textWidth, _ := pdf.MeasureTextWidth(text)
	xStart := xEnd - textWidth
	pdf.SetXY(xStart, y)
	pdf.Cell(nil, text)
}

func generatePDF(invoice Invoice) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Roboto", "./static/Roboto-Regular.ttf")
	if err != nil {
		return err
	}
	err = pdf.AddTTFFont("Roboto-Bold", "./static/Roboto-Bold.ttf")
	if err != nil {
		return err
	}

	rect := &gopdf.Rect{
		W: 167,
		H: 17,
	}
	err = pdf.Image("./static/hintermann-logo.png", 40, 32, rect)
	if err != nil {
		return err
	}
	rect = &gopdf.Rect{
		W: 138,
		H: 14,
	}
	err = pdf.Image("./static/hintermann-logo-small.png", 40, 796, rect)
	if err != nil {
		return err
	}

	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(94, 100, 112)

	// hintermann
	pdf.SetXY(40, 63)
	pdf.Cell(nil, "Asociația de Caritate Hintermann")
	pdf.SetXY(40, 79)
	pdf.Cell(nil, "Strada Spicului, Nr. 12")
	pdf.SetXY(40, 95)
	pdf.Cell(nil, "Bl. 40, Sc. A, Ap. 12")
	pdf.SetXY(40, 111)
	pdf.Cell(nil, "500460")
	pdf.SetXY(40, 127)
	pdf.Cell(nil, "Brașov")
	pdf.SetXY(40, 143)
	pdf.Cell(nil, "România")

	// cilent
	pdf.SetXY(312, 63)
	pdf.Cell(nil, "ID tranzacție:")
	rightAlign(&pdf, 555, 63, invoice.TransactionId)
	pdf.SetXY(312, 79)
	pdf.Cell(nil, "Data emiterii:")
	rightAlign(&pdf, 555, 79, invoice.IssueDate)
	pdf.SetXY(312, 95)
	pdf.Cell(nil, "Client:")
	rightAlign(&pdf, 555, 95, invoice.ClientName)

	// table head
	pdf.SetXY(40, 195)
	pdf.Cell(nil, "Serviciu")
	pdf.SetXY(312, 195)
	pdf.Cell(nil, "Cantitate")
	pdf.SetXY(419, 195)
	pdf.Cell(nil, "Preț unitar")
	pdf.SetXY(532, 195)
	pdf.Cell(nil, "Total")

	// product description
	pdf.SetXY(40, 253)
	pdf.Cell(nil, "Fiecare donație contribuie la transformarea")
	pdf.SetXY(40, 266)
	pdf.Cell(nil, "vieților familiilor românești aflate în mare nevoie.")
	pdf.SetXY(40, 279)
	pdf.Cell(nil, "Ia parte și tu acum.")

	// product price
	pdf.SetXY(347, 237)
	pdf.Cell(nil, "1")

	unitPrice := fmt.Sprintf("%.2f lei", invoice.UnitPrice)
	rightAlign(&pdf, 466, 237, unitPrice)

	total := fmt.Sprintf("%.2f lei", invoice.Total)
	rightAlign(&pdf, 555, 237, total)

	// footer
	pdf.SetXY(347, 796)
	pdf.Cell(nil, "contact@hintermann.ro")
	pdf.SetXY(492, 796)
	pdf.Cell(nil, "Pagina 1 din 1")

	// summary
	pdf.SetXY(312, 321)
	pdf.Cell(nil, "Subtotal:")
	pdf.SetXY(312, 343)
	pdf.Cell(nil, "TVA:")
	pdf.SetXY(312, 397)
	pdf.Cell(nil, "Debitat din plata dvs.:")
	pdf.SetXY(517, 321)
	pdf.Cell(nil, "10.00 lei")
	pdf.SetXY(522, 343)
	pdf.Cell(nil, "0.00 lei")
	pdf.SetXY(514, 397)
	pdf.Cell(nil, "-10.00 lei")

	// product title
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(40, 237)
	pdf.Cell(nil, invoice.ProductName)

	// summary bold
	pdf.SetFont("Roboto-Bold", "", 10)
	pdf.SetXY(312, 375)
	pdf.Cell(nil, "Total:")
	pdf.SetXY(515, 375)
	pdf.Cell(nil, "10.00 lei")
	pdf.SetXY(312, 429)
	pdf.Cell(nil, "Sumă datorată:")
	pdf.SetXY(521, 429)
	pdf.Cell(nil, "0.00 lei")

	// title
	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetXY(494, 32)
	pdf.Cell(nil, "Factură")

	// strokes
	pdf.SetStrokeColor(215, 218, 224)
	pdf.SetLineWidth(0.5)
	// table head
	pdf.Line(40, 216.5, 555, 216.5)
	// summary
	pdf.Line(40, 310, 555, 310)
	// summary total
	pdf.Line(312, 364.5, 555, 364.5)
	// summary amount due
	pdf.Line(312, 418.5, 555, 418.5)
	// footer
	pdf.Line(40, 775.5, 555, 775.5)

	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, "output.pdf")
	return pdf.WritePdf(pdfFile)
}

func main() {
	inv := Invoice{
		ClientName:    "Ungureanu Daniel",
		IssueDate:     "12 Aug, 2024",
		TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
		ProductName:   "Donație unică de 420 lei",
		UnitPrice:     420420.00,
		Total:         6969.00,
	}
	err := generatePDF(inv)
	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
