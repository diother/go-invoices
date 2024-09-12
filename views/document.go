package views

import (
	"github.com/diother/go-invoices/models"
	"github.com/signintech/gopdf"
)

type Document interface {
	GenerateDocument()
	addHeader()
	addFooter()
}

type InvoiceView struct {
	Invoice *models.Invoice
}

type PayoutView struct {
	Payout *models.Payout
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
