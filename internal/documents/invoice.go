package documents

import (
	"fmt"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/signintech/gopdf"
)

func (s *DocumentService) GenerateInvoice(donation *dto.FormattedDonation) (pdf *gopdf.GoPdf, err error) {
	pdf = &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err = setFonts(pdf); err != nil {
		return nil, fmt.Errorf("failed setting fonts: %w", err)
	}
	resetTextStyles(pdf)

	if err = addInvoiceHeader(pdf, donation); err != nil {
		return
	}
	if err = addInvoiceFooter(pdf); err != nil {
		return
	}

	addInvoiceTable(pdf)
	addInvoiceProduct(pdf, donation)
	addInvoiceSummary(pdf, donation)
	return
}

func addInvoiceHeader(pdf *gopdf.GoPdf, donation *dto.FormattedDonation) error {
	const startY = marginTop

	if err := addImage(pdf, "./static/images/hintermann-logo.png", marginLeft, marginTop, 167, 17); err != nil {
		return err
	}
	setText(pdf, marginLeft, startY+31, "Asociația de Caritate Hintermann")
	setText(pdf, marginLeft, startY+47, "Strada Spicului, Nr. 12")
	setText(pdf, marginLeft, startY+63, "Bl. 40, Sc. A, Ap. 12")
	setText(pdf, marginLeft, startY+79, "500460")
	setText(pdf, marginLeft, startY+95, "Brașov")
	setText(pdf, marginLeft, startY+111, "România")

	setText(pdf, 312, startY+31, "ID tranzacție:")
	setRightAlignedText(pdf, marginRight, startY+31, donation.ID)
	setText(pdf, 312, startY+47, "Data emiterii:")
	setRightAlignedText(pdf, marginRight, startY+47, donation.Created)
	setText(pdf, 312, startY+63, "Nume client:")
	setRightAlignedText(pdf, marginRight, startY+63, donation.ClientName)
	setText(pdf, 312, startY+79, "Email client:")
	setRightAlignedText(pdf, marginRight, startY+79, donation.ClientEmail)

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Factură")

	resetTextStyles(pdf)
	return nil
}

func addInvoiceFooter(pdf *gopdf.GoPdf) error {
	const endY = marginBottom

	if err := addImage(pdf, "./static/images/hintermann-logo-small.png", marginLeft, 796, 138, 14); err != nil {
		return fmt.Errorf("failed setting image: %w", err)
	}
	setRightAlignedText(pdf, 452, endY-14, "contact@hintermann.ro")
	setText(pdf, 492, endY-14, "Pagina 1 din 1")

	pdf.Line(marginLeft, endY-36.5, marginRight, endY-36.5)
	pdf.Line(471.5, endY-16, 471.5, endY-4)
	return nil
}

func addInvoiceTable(pdf *gopdf.GoPdf) {
	const startY = 195

	setText(pdf, marginLeft, startY, "Serviciu")
	setText(pdf, 312, startY, "Cantitate")
	setText(pdf, 419, startY, "Preț unitar")
	setText(pdf, 532, startY, "Total")

	pdf.Line(marginLeft, startY+21.5, marginRight, startY+21.5)
}

func addInvoiceProduct(pdf *gopdf.GoPdf, donation *dto.FormattedDonation) {
	const startY = 237

	setText(pdf, marginLeft, startY+16, "Fiecare donație contribuie la transformarea")
	setText(pdf, marginLeft, startY+29, "vieților familiilor românești aflate în mare nevoie.")
	setText(pdf, marginLeft, startY+42, "Ia parte și tu acum.")

	setText(pdf, 347, startY, "1")

	setRightAlignedText(pdf, 466, startY, donation.Gross)
	setRightAlignedText(pdf, marginRight, startY, donation.Gross)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY, "Donație de "+donation.Gross)
	pdf.SetTextColor(94, 100, 112)
}

func addInvoiceSummary(pdf *gopdf.GoPdf, donation *dto.FormattedDonation) {
	const startY = 311

	setText(pdf, 312, startY+10, "Subtotal:")
	setText(pdf, 312, startY+32, "TVA:")
	setText(pdf, 312, startY+86, "Debitat din plata dvs.:")

	setRightAlignedText(pdf, marginRight, startY+10, donation.Gross)

	setText(pdf, 522, startY+32, "0.00 lei")

	setRightAlignedText(pdf, marginRight, startY+86, "-"+donation.Gross)

	pdf.SetFont("Roboto-Bold", "", 10)
	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 312, startY+64, "Total:")

	setRightAlignedText(pdf, marginRight, startY+64, donation.Gross)

	setText(pdf, 312, startY+118, "Sumă datorată:")
	setText(pdf, 521, startY+118, "0.00 lei")

	pdf.Line(marginLeft, startY, marginRight, startY)
	pdf.Line(312, startY+53.5, marginRight, startY+53.5)
	pdf.Line(312, startY+107.5, marginRight, startY+107.5)

	resetTextStyles(pdf)
}

func setText(pdf *gopdf.GoPdf, x, y float64, text string) {
	pdf.SetXY(x, y)
	pdf.Cell(nil, text)
}

func setRightAlignedText(pdf *gopdf.GoPdf, xEnd, y float64, text string) {
	textWidth, _ := pdf.MeasureTextWidth(text)
	xStart := xEnd - textWidth
	setText(pdf, xStart, y, text)
}

func addImage(pdf *gopdf.GoPdf, path string, x, y, w, h float64) error {
	rect := &gopdf.Rect{W: w, H: h}
	return pdf.Image(path, x, y, rect)
}

func resetTextStyles(pdf *gopdf.GoPdf) {
	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(94, 100, 112)

	pdf.SetStrokeColor(215, 218, 224)
	pdf.SetLineWidth(0.5)
}

func setFonts(pdf *gopdf.GoPdf) error {
	if err := pdf.AddTTFFont("Roboto", "./static/fonts/Roboto-Regular.ttf"); err != nil {
		return err
	}
	return pdf.AddTTFFont("Roboto-Bold", "./static/fonts/Roboto-Bold.ttf")
}
