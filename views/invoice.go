package views

import (
	"fmt"
	"path/filepath"

	"github.com/diother/go-invoices/models"
	"github.com/signintech/gopdf"
)

func addFooter(pdf *gopdf.GoPdf) {
	setText(pdf, 347, 796, "contact@hintermann.ro")
	setText(pdf, 492, 796, "Pagina 1 din 1")
}

func addHeader(pdf *gopdf.GoPdf, invoice *models.Invoice) {
	setText(pdf, 40, 63, "Asociația de Caritate Hintermann")
	setText(pdf, 40, 79, "Strada Spicului, Nr. 12")
	setText(pdf, 40, 95, "Bl. 40, Sc. A, Ap. 12")
	setText(pdf, 40, 111, "500460")
	setText(pdf, 40, 127, "Brașov")
	setText(pdf, 40, 143, "România")

	setText(pdf, 312, 63, "ID tranzacție:")
	setRightAlignedText(pdf, 555, 63, invoice.TransactionId)
	setText(pdf, 312, 79, "Data emiterii:")
	setRightAlignedText(pdf, 555, 79, invoice.IssueDate)
	setText(pdf, 312, 95, "Client:")
	setRightAlignedText(pdf, 555, 95, invoice.ClientName)

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 494, 32, "Factură")

	resetTextStyles(pdf)
}

func addTable(pdf *gopdf.GoPdf) {
	setText(pdf, 40, 195, "Serviciu")
	setText(pdf, 312, 195, "Cantitate")
	setText(pdf, 419, 195, "Preț unitar")
	setText(pdf, 532, 195, "Total")
}

func addSummary(pdf *gopdf.GoPdf, invoice *models.Invoice) {
	setText(pdf, 312, 321, "Subtotal:")
	setText(pdf, 312, 343, "TVA:")
	setText(pdf, 312, 397, "Debitat din plata dvs.:")

	total := fmt.Sprintf("%.2f lei", invoice.Total)
	setRightAlignedText(pdf, 555, 321, total)

	setText(pdf, 522, 343, "0.00 lei")

	minusTotal := fmt.Sprintf("-%.2f lei", invoice.Total)
	setRightAlignedText(pdf, 555, 397, minusTotal)

	pdf.SetFont("Roboto-Bold", "", 10)
	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 312, 375, "Total:")

	setRightAlignedText(pdf, 555, 375, total)

	setText(pdf, 312, 429, "Sumă datorată:")
	setText(pdf, 521, 429, "0.00 lei")

	resetTextStyles(pdf)
}

func addProduct(pdf *gopdf.GoPdf, invoice *models.Invoice) {
	setText(pdf, 40, 253, "Fiecare donație contribuie la transformarea")
	setText(pdf, 40, 266, "vieților familiilor românești aflate în mare nevoie.")
	setText(pdf, 40, 279, "Ia parte și tu acum.")

	setText(pdf, 347, 237, "1")

	unitPrice := fmt.Sprintf("%.2f lei", invoice.UnitPrice)
	setRightAlignedText(pdf, 466, 237, unitPrice)

	total := fmt.Sprintf("%.2f lei", invoice.Total)
	setRightAlignedText(pdf, 555, 237, total)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 40, 237, invoice.ProductName)
	pdf.SetTextColor(94, 100, 112)
}

func addStrokes(pdf *gopdf.GoPdf) {
	pdf.SetStrokeColor(215, 218, 224)
	pdf.SetLineWidth(0.5)

	pdf.Line(40, 216.5, 555, 216.5)
	pdf.Line(40, 310, 555, 310)
	pdf.Line(312, 364.5, 555, 364.5)
	pdf.Line(312, 418.5, 555, 418.5)
	pdf.Line(40, 775.5, 555, 775.5)
}

func GenerateInvoicePdf(invoice *models.Invoice) error {
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

	err = addImage(&pdf, "./static/hintermann-logo.png", 40, 32, 167, 17)
	if err != nil {
		return err
	}
	err = addImage(&pdf, "./static/hintermann-logo-small.png", 40, 796, 138, 14)
	if err != nil {
		return err
	}

	resetTextStyles(&pdf)

	addHeader(&pdf, invoice)
	addFooter(&pdf)
	addTable(&pdf)
	addProduct(&pdf, invoice)
	addSummary(&pdf, invoice)
	addStrokes(&pdf)

	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, "output.pdf")
	return pdf.WritePdf(pdfFile)
}
