package documents

import (
	"fmt"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/signintech/gopdf"
)

func (s DocumentService) GeneratePayoutReport(payoutReportData *dto.PayoutReportData) (pdf *gopdf.GoPdf, err error) {
	payout := payoutReportData.Payout
	items := payoutReportData.Items

	pdf = &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err = setFonts(pdf); err != nil {
		return nil, fmt.Errorf("failed setting fonts: %w", err)
	}
	resetTextStyles(pdf)

	itemsLength := len(items)
	pagesNeeded := pagesNeeded(itemsLength)
	currentPage := 1

	if err = addPayoutReportHeader(pdf, payout.Created); err != nil {
		return nil, fmt.Errorf("failed adding the header: %w", err)
	}
	if err = addPayoutReportFooter(pdf, currentPage, pagesNeeded); err != nil {
		return nil, fmt.Errorf("failed adding the footer: %w", err)
	}

	addPayoutSummary(pdf, payout)
	addPayoutTable(pdf, firstPageTableY)

	currentY := firstPageStartY
	maxItemsPerPage := firstPageCapacity

	var itemCounter int
	for _, item := range items {
		if itemCounter == maxItemsPerPage {
			pdf.AddPage()
			currentPage++

			if err = addPayoutReportSecondaryHeader(pdf); err != nil {
				return nil, fmt.Errorf("failed adding the secondary footer: %w", err)
			}
			if err = addPayoutReportFooter(pdf, currentPage, pagesNeeded); err != nil {
				return nil, fmt.Errorf("failed adding the footer: %w", err)
			}

			addPayoutTable(pdf, subsequentPageTableY)

			currentY = secondPageStartY
			itemCounter = 0
			maxItemsPerPage = subsequentPageCapacity
		}
		addPayoutProduct(pdf, item, currentY)
		currentY += itemHeight
		itemCounter++
	}
	return
}

func addPayoutReportHeader(pdf *gopdf.GoPdf, created string) error {
	const startY = marginTop

	if err := addImage(pdf, "./static/images/stripe-logo.png", marginLeft, startY, 51, 21); err != nil {
		return err
	}

	setText(pdf, marginLeft, startY+31, "Stripe Payments Europe, Limited")
	setText(pdf, marginLeft, startY+47, "The One Building")
	setText(pdf, marginLeft, startY+63, "1 Grand Canal Street Lower")
	setText(pdf, marginLeft, startY+79, "Dublin 2")
	setText(pdf, marginLeft, startY+95, "Co. Dublin")
	setText(pdf, marginLeft, startY+111, "Ireland")

	setText(pdf, 312, startY+31, "Data emiterii:")
	setRightAlignedText(pdf, marginRight, startY+31, created)
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
	return nil
}

func addPayoutReportSecondaryHeader(pdf *gopdf.GoPdf) error {
	const startY = marginTop

	if err := addImage(pdf, "./static/images/stripe-logo.png", marginLeft, startY, 51, 21); err != nil {
		return err
	}
	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras plată")

	resetTextStyles(pdf)
	return nil
}

func addPayoutReportFooter(pdf *gopdf.GoPdf, currentPage int, pagesNeeded int) error {
	const endY = marginBottom
	if err := addImage(pdf, "./static/images/stripe-logo-small.png", marginLeft, endY-17, 41, 17); err != nil {
		return err
	}

	pageInfo := fmt.Sprintf("Pagina %d din %d", currentPage, pagesNeeded)
	setText(pdf, 492, endY-15.5, pageInfo)

	pdf.Line(marginLeft, endY-37, marginRight, endY-37)
	return nil
}

func addPayoutSummary(pdf *gopdf.GoPdf, payout *dto.FormattedPayout) {
	const startY = 211

	setText(pdf, 81, startY+10, payout.ID)
	setText(pdf, 112, startY+26, payout.Created)

	setText(pdf, 312, startY+10, "Preț brut:")
	setText(pdf, 312, startY+26, "Taxe Stripe:")

	setRightAlignedText(pdf, marginRight, startY+10, payout.Gross)
	setRightAlignedText(pdf, marginRight, startY+26, "-"+payout.Fee)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY+10, "ID plată:")
	setText(pdf, marginLeft, startY+26, "Data efectuării:")

	pdf.SetFont("Roboto-Bold", "", 10)
	setText(pdf, 312, startY+42, "Total:")
	setRightAlignedText(pdf, marginRight, startY+42, payout.Net)

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

func addPayoutProduct(pdf *gopdf.GoPdf, item *dto.PayoutReportItem, startY float64) {
	setText(pdf, marginLeft, startY+16, item.ID)

	setRightAlignedText(pdf, 367, startY, item.Gross)
	setRightAlignedText(pdf, 474, startY, "-"+item.Fee)

	var productName string
	if item.Type == "donation" {
		setRightAlignedText(pdf, marginRight, startY, item.Net)
		pdf.SetTextColor(0, 0, 0)
		productName = "Donație de " + item.Gross
	} else {
		setRightAlignedText(pdf, marginRight, startY, "-"+item.Net)
		pdf.SetTextColor(0, 0, 0)
		productName = "Taxă Stripe de " + item.Fee
	}
	setText(pdf, marginLeft, startY, productName)
	pdf.SetTextColor(94, 100, 112)
}

// needs unit test
func pagesNeeded(itemsLength int) int {
	const (
		firstPageCapacity      = 8
		subsequentPageCapacity = 12
	)
	remainingItems := itemsLength - firstPageCapacity
	var totalPages int

	if remainingItems > 0 {
		additionalPages := (remainingItems + subsequentPageCapacity - 1) / subsequentPageCapacity
		totalPages = 1 + additionalPages
	} else {
		totalPages = 1
	}
	return totalPages
}
