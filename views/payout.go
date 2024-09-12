package views

import (
	"fmt"
	"path/filepath"

	"github.com/diother/go-invoices/models"
	"github.com/signintech/gopdf"
)

func (p PayoutView) addHeader(pdf *gopdf.GoPdf, payout *models.Payout) {
	setText(pdf, 40, 63, "Stripe Payments Europe, Limited")
	setText(pdf, 40, 79, "The One Building")
	setText(pdf, 40, 95, "1 Grand Canal Street Lower")
	setText(pdf, 40, 111, "Dublin 2")
	setText(pdf, 40, 127, "Co. Dublin")
	setText(pdf, 40, 143, "Ireland")

	setText(pdf, 312, 63, "Data emiterii:")
	setRightAlignedText(pdf, 555, 63, payout.IssueDate)
	setText(pdf, 312, 79, "Nr. cont:")
	setRightAlignedText(pdf, 555, 79, "acct_1PVfUvDXCtuWOFq8")
	setText(pdf, 312, 95, "Proprietar cont:")
	setRightAlignedText(pdf, 555, 95, "Asociația de Caritate Hintermann")
	setText(pdf, 312, 111, "Adresă:")
	setRightAlignedText(pdf, 555, 111, "Strada Spicului, Nr. 12")
	setRightAlignedText(pdf, 555, 127, "Bl. 40, Sc. A, Ap. 12")
	setRightAlignedText(pdf, 555, 143, "Brașov, România")
	setRightAlignedText(pdf, 555, 159, "500460")

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, 555, 32, "Extras plată")

	resetTextStyles(pdf)
}

func (p PayoutView) addFooter(pdf *gopdf.GoPdf) {
	setText(pdf, 347, 794.5, "support@stripe.com")
	setText(pdf, 492, 794.5, "Pagina 1 din 1")
}

func addPayoutTable(pdf *gopdf.GoPdf) {
	setText(pdf, 40, 315, "Tranzacție")
	setText(pdf, 328, 315, "Preț brut")
	setText(pdf, 424.5, 315, "Taxă Stripe")
	setText(pdf, 532, 315, "Total")
}

func addPayoutSummary(pdf *gopdf.GoPdf, payout *models.Payout) {
	setText(pdf, 81, 221, "po_1PZ0rmDXCtuWOFq8n33WSnN9")
	setText(pdf, 112, 237, "31 Jul, 2024")

	setText(pdf, 312, 221, "Preț brut:")
	setText(pdf, 312, 237, "Taxe Stripe:")

	setRightAlignedText(pdf, 555, 221, "10 lei")
	setRightAlignedText(pdf, 555, 237, "10 lei")

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 40, 221, "ID plată:")
	setText(pdf, 40, 237, "Data efectuării:")

	pdf.SetFont("Roboto-Bold", "", 10)
	setText(pdf, 312, 253, "Total:")
	setRightAlignedText(pdf, 555, 253, "10 lei")

	resetTextStyles(pdf)
}

func addPayoutProduct(pdf *gopdf.GoPdf, payout *models.Payout) {
	setText(pdf, 40, 373, "ch_3PXmP9DXCtuWOFq82TlDyZ3F")

	setRightAlignedText(pdf, 367, 357, "10.00 lei")

	unitPrice := fmt.Sprintf("%.2f lei", payout.UnitPrice)
	setRightAlignedText(pdf, 474, 357, unitPrice)

	total := fmt.Sprintf("%.2f lei", payout.Total)
	setRightAlignedText(pdf, 555, 357, total)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 40, 357, payout.ProductName)
	pdf.SetTextColor(94, 100, 112)
}

func (p PayoutView) GenerateDocument(payout *models.Payout) error {
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

	err = addImage(&pdf, "./static/stripe-logo.png", 40, 32, 51, 21)
	if err != nil {
		return err
	}
	err = addImage(&pdf, "./static/stripe-logo-small.png", 40, 793, 41, 17)
	if err != nil {
		return err
	}

	resetTextStyles(&pdf)

	p.addHeader(&pdf, payout)
	p.addFooter(&pdf)

	addPayoutTable(&pdf)
	addPayoutProduct(&pdf, payout)
	addPayoutSummary(&pdf, payout)

	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, "output.pdf")
	return pdf.WritePdf(pdfFile)
}
