package views

import (
	"fmt"
	"path/filepath"

	"github.com/diother/go-invoices/models"
	"github.com/signintech/gopdf"
)

const (
	marginTop   = 32
	marginLeft  = 40
	marginRight = 555
)

func (p PayoutView) addHeader(pdf *gopdf.GoPdf, payout *models.Payout) {
	const startY = marginTop

	setText(pdf, marginLeft, startY+31, "Stripe Payments Europe, Limited")
	setText(pdf, marginLeft, startY+47, "The One Building")
	setText(pdf, marginLeft, startY+63, "1 Grand Canal Street Lower")
	setText(pdf, marginLeft, startY+79, "Dublin 2")
	setText(pdf, marginLeft, startY+95, "Co. Dublin")
	setText(pdf, marginLeft, startY+111, "Ireland")

	setText(pdf, 312, startY+31, "Data emiterii:")
	setRightAlignedText(pdf, marginRight, startY+31, payout.IssueDate)
	setText(pdf, 312, startY+47, "Nr. cont:")
	setRightAlignedText(pdf, marginRight, startY+47, "acct_1PVfUvDXCtuWOFq8")
	setText(pdf, 312, startY+63, "Proprietar cont:")
	setRightAlignedText(pdf, marginRight, startY+63, "Asociația de Caritate Hintermann")
	setText(pdf, 312, startY+79, "Adresă:")
	setRightAlignedText(pdf, marginRight, startY+79, "Strada Spicului, Nr. 12")
	setRightAlignedText(pdf, marginRight, startY+95, "Bl. MarginLeft, Sc. A, Ap. 12")
	setRightAlignedText(pdf, marginRight, startY+111, "Brașov, România")
	setRightAlignedText(pdf, marginRight, startY+127, "500460")

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras plată")

	resetTextStyles(pdf)
}

func (p PayoutView) addFooter(pdf *gopdf.GoPdf) {
	setText(pdf, 347, 794.5, "support@stripe.com")
	setText(pdf, 492, 794.5, "Pagina 1 din 1")
}

func addPayoutTable(pdf *gopdf.GoPdf) {
	setText(pdf, marginLeft, 315, "Tranzacție")
	setText(pdf, 328, 315, "Preț brut")
	setText(pdf, 424.5, 315, "Taxă Stripe")
	setText(pdf, 532, 315, "Total")
}

func addPayoutSummary(pdf *gopdf.GoPdf, payout *models.Payout) {
	setText(pdf, 81, 221, "po_1PZ0rmDXCtuWOFq8n33WSnN9")
	setText(pdf, 112, 237, "31 Jul, 2024")

	setText(pdf, 312, 221, "Preț brut:")
	setText(pdf, 312, 237, "Taxe Stripe:")

	setRightAlignedText(pdf, marginRight, 221, "10 lei")
	setRightAlignedText(pdf, marginRight, 237, "10 lei")

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, 221, "ID plată:")
	setText(pdf, marginLeft, 237, "Data efectuării:")

	pdf.SetFont("Roboto-Bold", "", 10)
	setText(pdf, 312, 253, "Total:")
	setRightAlignedText(pdf, marginRight, 253, "10 lei")

	resetTextStyles(pdf)
}

func addPayoutProduct(pdf *gopdf.GoPdf, payout *models.Payout) {
	setText(pdf, marginLeft, 373, "ch_3PXmP9DXCtuWOFq82TlDyZ3F")

	setRightAlignedText(pdf, 367, 357, "10.00 lei")

	unitPrice := fmt.Sprintf("%.2f lei", payout.UnitPrice)
	setRightAlignedText(pdf, 474, 357, unitPrice)

	total := fmt.Sprintf("%.2f lei", payout.Total)
	setRightAlignedText(pdf, marginRight, 357, total)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, 357, payout.ProductName)
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

	err = addImage(&pdf, "./static/stripe-logo.png", marginLeft, 32, 51, 21)
	if err != nil {
		return err
	}
	err = addImage(&pdf, "./static/stripe-logo-small.png", marginLeft, 793, 41, 17)
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
