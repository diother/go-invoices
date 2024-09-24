package documents

import (
	"fmt"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/signintech/gopdf"
)

func (s DocumentService) GenerateMonthlyReport(payouts []*dto.FormattedPayout) (pdf *gopdf.GoPdf, err error) {
	pdf = &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	pdf = &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err = setFonts(pdf); err != nil {
		return nil, fmt.Errorf("failed setting fonts: %w", err)
	}
	resetTextStyles(pdf)

	itemsLength := len(payouts)
	pagesNeeded := pagesNeeded(itemsLength)
	currentPage := 1

	if err = addMonthlyReportHeader(pdf, "placeholder date"); err != nil {
		return nil, fmt.Errorf("failed adding the header: %w", err)
	}
	if err = addMonthlyReportFooter(pdf, currentPage, pagesNeeded); err != nil {
		return nil, fmt.Errorf("failed adding the footer: %w", err)
	}

	addMonthlyPayoutSummary(pdf)
	addMonthlyPayoutTable(pdf, firstPageTableY)

	currentY := firstPageStartY
	maxItemsPerPage := firstPageCapacity

	var itemCounter int
	for _, payout := range payouts {
		if itemCounter == maxItemsPerPage {
			pdf.AddPage()
			currentPage++

			if err = addMonthlyReportSecondaryHeader(pdf); err != nil {
				return nil, fmt.Errorf("failed adding the secondary header: %w", err)
			}
			if err = addMonthlyReportFooter(pdf, currentPage, pagesNeeded); err != nil {
				return nil, fmt.Errorf("failed adding the footer: %w", err)
			}

			addMonthlyPayoutTable(pdf, subsequentPageTableY)

			currentY = secondPageStartY
			itemCounter = 0
			maxItemsPerPage = subsequentPageCapacity
		}
		addMonthlyPayoutProduct(pdf, payout, currentY)
		currentY += itemHeight
		itemCounter++
	}
	return
}

func addMonthlyReportHeader(pdf *gopdf.GoPdf, created string) error {
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
	setRightAlignedText(pdf, marginRight, startY, "Extras lunar")

	resetTextStyles(pdf)
	return nil
}

func addMonthlyReportSecondaryHeader(pdf *gopdf.GoPdf) error {
	const startY = marginTop

	if err := addImage(pdf, "./static/images/stripe-logo.png", marginLeft, startY, 51, 21); err != nil {
		return err
	}
	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Extras lunar")

	resetTextStyles(pdf)
	return nil
}

func addMonthlyReportFooter(pdf *gopdf.GoPdf, currentPage, pagesNeeded int) error {
	const endY = marginBottom

	if err := addImage(pdf, "./static/images/stripe-logo-small.png", marginLeft, endY-17, 41, 17); err != nil {
		return err
	}
	pageInfo := fmt.Sprintf("Pagina %d din %d", currentPage, pagesNeeded)
	setText(pdf, 492, endY-15.5, pageInfo)

	pdf.Line(marginLeft, endY-37, marginRight, endY-37)
	return nil
}

func addMonthlyPayoutSummary(pdf *gopdf.GoPdf) {
	const startY = 211

	setText(pdf, marginLeft, startY+26, "placeholder period")

	setText(pdf, 312, startY+10, "Preț brut:")
	setText(pdf, 312, startY+26, "Taxe Stripe:")

	setRightAlignedText(pdf, marginRight, startY+10, "placeholder lei")
	setRightAlignedText(pdf, marginRight, startY+26, "placeholder lei")

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY+10, "Periodă extras:")

	pdf.SetFont("Roboto-Bold", "", 10)
	setText(pdf, 312, startY+42, "Total:")
	setRightAlignedText(pdf, marginRight, startY+42, "placeholder lei")

	resetTextStyles(pdf)

	pdf.Line(marginLeft, startY-.5, marginRight, startY-.5)
	pdf.Line(marginLeft, startY+63.5, marginRight, startY+63.5)
	pdf.Line(297.5, startY-.5, 298.5, startY+63.5)
}

func addMonthlyPayoutTable(pdf *gopdf.GoPdf, startY float64) {
	setText(pdf, marginLeft, startY, "Tranzacție")
	setText(pdf, 328, startY, "Preț brut")
	setText(pdf, 424.5, startY, "Taxă Stripe")
	setText(pdf, 532, startY, "Total")

	pdf.Line(marginLeft, startY+21.5, marginRight, startY+21.5)
}

func addMonthlyPayoutProduct(pdf *gopdf.GoPdf, payout *dto.FormattedPayout, startY float64) {
	setText(pdf, marginLeft, startY+16, payout.Created)

	setRightAlignedText(pdf, 367, startY, payout.Gross)
	setRightAlignedText(pdf, 474, startY, "-"+payout.Fee)
	setRightAlignedText(pdf, marginRight, startY, payout.Net)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY, payout.ID)
	pdf.SetTextColor(94, 100, 112)
}
