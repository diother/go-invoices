package pdf

import (
	"path/filepath"

	"github.com/diother/go-invoices/internal/models"
	"github.com/signintech/gopdf"
)

func (i InvoicePdf) addHeader(pdf *gopdf.GoPdf, invoice *models.Invoice) error {
	const startY = marginTop

	err := addImage(pdf, "./internal/pdf/static/images/hintermann-logo.png", marginLeft, marginTop, 167, 17)
	if err != nil {
		return err
	}

	setText(pdf, marginLeft, startY+31, "Asociația de Caritate Hintermann")
	setText(pdf, marginLeft, startY+47, "Strada Spicului, Nr. 12")
	setText(pdf, marginLeft, startY+63, "Bl. marginLeft, Sc. A, Ap. 12")
	setText(pdf, marginLeft, startY+79, "500460")
	setText(pdf, marginLeft, startY+95, "Brașov")
	setText(pdf, marginLeft, startY+111, "România")

	setText(pdf, 312, startY+31, "ID tranzacție:")
	setRightAlignedText(pdf, marginRight, startY+31, invoice.TransactionId)
	setText(pdf, 312, startY+47, "Data emiterii:")
	setRightAlignedText(pdf, marginRight, startY+47, invoice.IssueDate)
	setText(pdf, 312, startY+63, "Client:")
	setRightAlignedText(pdf, marginRight, startY+63, invoice.ClientName)

	pdf.SetFont("Roboto-Bold", "", 18)
	pdf.SetTextColor(0, 0, 0)
	setRightAlignedText(pdf, marginRight, startY, "Factură")

	resetTextStyles(pdf)
	return nil
}

func (i InvoicePdf) addFooter(pdf *gopdf.GoPdf) error {
	const endY = marginBottom

	err := addImage(pdf, "./internal/pdf/static/images/hintermann-logo-small.png", marginLeft, 796, 138, 14)
	if err != nil {
		return err
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

func addInvoiceProduct(pdf *gopdf.GoPdf, invoice *models.Invoice) {
	const startY = 237

	setText(pdf, marginLeft, startY+16, "Fiecare donație contribuie la transformarea")
	setText(pdf, marginLeft, startY+29, "vieților familiilor românești aflate în mare nevoie.")
	setText(pdf, marginLeft, startY+42, "Ia parte și tu acum.")

	setText(pdf, 347, startY, "1")

	setRightAlignedText(pdf, 466, startY, invoice.UnitPrice)
	setRightAlignedText(pdf, marginRight, startY, invoice.Total)

	pdf.SetTextColor(0, 0, 0)
	setText(pdf, marginLeft, startY, invoice.ProductName)
	pdf.SetTextColor(94, 100, 112)
}

func addInvoiceSummary(pdf *gopdf.GoPdf, invoice *models.Invoice) {
	const startY = 311

	setText(pdf, 312, startY+10, "Subtotal:")
	setText(pdf, 312, startY+32, "TVA:")
	setText(pdf, 312, startY+86, "Debitat din plata dvs.:")

	setRightAlignedText(pdf, marginRight, startY+10, invoice.Total)

	setText(pdf, 522, startY+32, "0.00 lei")

	setRightAlignedText(pdf, marginRight, startY+86, "-"+invoice.Total)

	pdf.SetFont("Roboto-Bold", "", 10)
	pdf.SetTextColor(0, 0, 0)
	setText(pdf, 312, startY+64, "Total:")

	setRightAlignedText(pdf, marginRight, startY+64, invoice.Total)

	setText(pdf, 312, startY+118, "Sumă datorată:")
	setText(pdf, 521, startY+118, "0.00 lei")

	pdf.Line(marginLeft, startY, marginRight, startY)
	pdf.Line(312, startY+53.5, marginRight, startY+53.5)
	pdf.Line(312, startY+107.5, marginRight, startY+107.5)

	resetTextStyles(pdf)
}

func (i InvoicePdf) GenerateDocument() error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Roboto", "./internal/pdf/static/fonts/Roboto-Regular.ttf")
	if err != nil {
		return err
	}
	err = pdf.AddTTFFont("Roboto-Bold", "./internal/pdf/static/fonts/Roboto-Bold.ttf")
	if err != nil {
		return err
	}

	resetTextStyles(&pdf)

	err = i.addHeader(&pdf, i.Invoice)
	if err != nil {
		return err
	}
	err = i.addFooter(&pdf)
	if err != nil {
		return err
	}

	addInvoiceTable(&pdf)
	addInvoiceProduct(&pdf, i.Invoice)
	addInvoiceSummary(&pdf, i.Invoice)

	outputDir := "./internal/pdf/output/"
	pdfFile := filepath.Join(outputDir, "output.pdf")
	return pdf.WritePdf(pdfFile)
}
