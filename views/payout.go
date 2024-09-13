package views

import (
	"fmt"
	"path/filepath"

	"github.com/diother/go-invoices/models"
	"github.com/signintech/gopdf"
)

func (p PayoutView) addHeader(pdf *gopdf.GoPdf, payout *models.Payout) {
	const startY = marginTop

	addImage(pdf, "./static/stripe-logo.png", marginLeft, startY, 51, 21)

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
	setRightAlignedText(pdf, marginRight, startY+95, "Bl. 40, Sc. A, Ap. 12")
	setRightAlignedText(pdf, marginRight, startY+111, "Brașov, România")
	setRightAlignedText(pdf, marginRight, startY+127, "500460")

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras plată")

	resetTextStyles(pdf)
}

func (p PayoutView) addSecondaryHeader(pdf *gopdf.GoPdf) {
	const startY = marginTop

	addImage(pdf, "./static/stripe-logo.png", marginLeft, startY, 51, 21)

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras plată")

	resetTextStyles(pdf)
}

func (p PayoutView) addFooter(pdf *gopdf.GoPdf) {
	const endY = marginBottom

	addImage(pdf, "./static/stripe-logo-small.png", marginLeft, endY-17, 41, 17)
	setText(pdf, 492, endY-15.5, "Pagina 1 din 1")
	pdf.Line(marginLeft, endY-37, marginRight, endY-37)
}

func addPayoutSummary(pdf *gopdf.GoPdf, payout *models.Payout) {
	const startY = 211

	setText(pdf, 81, startY+10, "po_1PZ0rmDXCtuWOFq8n33WSnN9")
	setText(pdf, 112, startY+26, "31 Jul, 2024")

	setText(pdf, 312, startY+10, "Preț brut:")
	setText(pdf, 312, startY+26, "Taxe Stripe:")

	setRightAlignedText(pdf, marginRight, startY+10, "10 lei")
	setRightAlignedText(pdf, marginRight, startY+26, "10 lei")

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY+10, "ID plată:")
	setText(pdf, marginLeft, startY+26, "Data efectuării:")

	pdf.SetFont("Roboto-Bold", "", 10)
	setText(pdf, 312, startY+42, "Total:")
	setRightAlignedText(pdf, marginRight, startY+42, "10 lei")

	resetTextStyles(pdf)

	pdf.Line(marginLeft, startY-.5, marginRight, startY-.5)
	pdf.Line(marginLeft, startY+63.5, marginRight, startY+63.5)
	pdf.Line(297.5, startY-.5, 298.5, startY+63.5)
}

func addPayoutTable(pdf *gopdf.GoPdf, startY float64) {
	setText(pdf, marginLeft, startY, "Tranzacție")
	setText(pdf, 328, startY, "Preț brut")
	setText(pdf, 424.5, startY, "Taxă Stripe")
	setText(pdf, 532, startY, "Total")

	pdf.Line(marginLeft, startY+21.5, marginRight, startY+21.5)
}

func addPayoutProduct(pdf *gopdf.GoPdf, payout *models.Payout, startY float64) {
	setText(pdf, marginLeft, startY+16, "ch_3PXmP9DXCtuWOFq82TlDyZ3F")

	setRightAlignedText(pdf, 367, startY, "10.00 lei")

	unitPrice := fmt.Sprintf("%.2f lei", payout.UnitPrice)
	setRightAlignedText(pdf, 474, startY, unitPrice)

	total := fmt.Sprintf("%.2f lei", payout.Total)
	setRightAlignedText(pdf, marginRight, startY, total)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY, payout.ProductName)
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

	const (
		itemHeight         = 50
		firstPageStartY    = 357.0
		secondPageStartY   = 135.0
		maxItemsFirstPage  = 8
		maxItemsSecondPage = 12
		firstPageTableY    = 315
		secondPageTableY   = 93
	)

	resetTextStyles(&pdf)

	p.addHeader(&pdf, payout)
	p.addFooter(&pdf)
	addPayoutSummary(&pdf, payout)
	addPayoutTable(&pdf, firstPageTableY)

	currentY := firstPageStartY
	maxItemsPerPage := maxItemsFirstPage

	var itemCounter int
	for i := 0; i < 45; i++ {
		if itemCounter == maxItemsPerPage {
			pdf.AddPage()

			p.addSecondaryHeader(&pdf)
			p.addFooter(&pdf)

			addPayoutTable(&pdf, secondPageTableY)

			currentY = secondPageStartY
			itemCounter = 0
			maxItemsPerPage = maxItemsSecondPage
		}
		addPayoutProduct(&pdf, payout, currentY)
		currentY += itemHeight
		itemCounter++
	}

	outputDir := "./pdf"
	pdfFile := filepath.Join(outputDir, "output.pdf")
	return pdf.WritePdf(pdfFile)
}
