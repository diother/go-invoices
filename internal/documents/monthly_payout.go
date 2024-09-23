package documents

// func (s DocumentService) addHeader(pdf *gopdf.GoPdf, issueDate string) error {
// 	const startY = marginTop
//
// 	err := addImage(pdf, "./internal/pdf/static/images/stripe-logo.png", marginLeft, startY, 51, 21)
// 	if err != nil {
// 		return err
// 	}
//
// 	setText(pdf, marginLeft, startY+31, "Stripe Payments Europe, Limited")
// 	setText(pdf, marginLeft, startY+47, "The One Building")
// 	setText(pdf, marginLeft, startY+63, "1 Grand Canal Street Lower")
// 	setText(pdf, marginLeft, startY+79, "Dublin 2")
// 	setText(pdf, marginLeft, startY+95, "Co. Dublin")
// 	setText(pdf, marginLeft, startY+111, "Ireland")
//
// 	setText(pdf, 312, startY+31, "Data emiterii:")
// 	setRightAlignedText(pdf, marginRight, startY+31, issueDate)
// 	setText(pdf, 312, startY+47, "Nr. cont:")
// 	setRightAlignedText(pdf, marginRight, startY+47, "acct_1PVfUvDXCtuWOFq8")
// 	setText(pdf, 312, startY+63, "Proprietar cont:")
// 	setRightAlignedText(pdf, marginRight, startY+63, "Asociația de Caritate Hintermann")
// 	setText(pdf, 312, startY+79, "Adresă:")
// 	setRightAlignedText(pdf, marginRight, startY+79, "Strada Spicului, Nr. 12")
// 	setRightAlignedText(pdf, marginRight, startY+95, "Bl. 40, Sc. A, Ap. 12")
// 	setRightAlignedText(pdf, marginRight, startY+111, "Brașov, România")
// 	setRightAlignedText(pdf, marginRight, startY+127, "500460")
//
// 	pdf.SetFont("Roboto-Bold", "", 18)
// 	pdf.SetTextColor(0, 0, 0)
// 	setRightAlignedText(pdf, marginRight, startY, "Extras lunar")
//
// 	resetTextStyles(pdf)
// 	return nil
// }
//
// func (s DocumentService) addSecondaryHeader(pdf *gopdf.GoPdf) error {
// 	const startY = marginTop
//
// 	err := addImage(pdf, "./internal/pdf/static/images/stripe-logo.png", marginLeft, startY, 51, 21)
// 	if err != nil {
// 		return err
// 	}
//
// 	pdf.SetFont("Roboto-Bold", "", 18)
// 	pdf.SetTextColor(0, 0, 0)
// 	setRightAlignedText(pdf, marginRight, startY, "Extras lunar")
//
// 	resetTextStyles(pdf)
// 	return nil
// }
//
// func (s DocumentService) addFooter(pdf *gopdf.GoPdf, pagesNeeded int, currentPage int) error {
// 	const endY = marginBottom
//
// 	err := addImage(pdf, "./internal/pdf/static/images/stripe-logo-small.png", marginLeft, endY-17, 41, 17)
// 	if err != nil {
// 		return err
// 	}
//
// 	pageInfo := fmt.Sprintf("Pagina %d din %d", currentPage, pagesNeeded)
// 	setText(pdf, 492, endY-15.5, pageInfo)
//
// 	pdf.Line(marginLeft, endY-37, marginRight, endY-37)
// 	return nil
// }
//
// func addMonthlyPayoutSummary(pdf *gopdf.GoPdf, payout *models.MonthlyPayout) {
// 	const startY = 211
//
// 	setText(pdf, marginLeft, startY+26, payout.ReportPeriod)
//
// 	setText(pdf, 312, startY+10, "Preț brut:")
// 	setText(pdf, 312, startY+26, "Taxe Stripe:")
//
// 	setRightAlignedText(pdf, marginRight, startY+10, payout.Gross)
// 	setRightAlignedText(pdf, marginRight, startY+26, payout.StripeFees)
//
// 	pdf.SetTextColor(0, 0, 0)
// 	setText(pdf, marginLeft, startY+10, "Periodă extras:")
//
// 	pdf.SetFont("Roboto-Bold", "", 10)
// 	setText(pdf, 312, startY+42, "Total:")
// 	setRightAlignedText(pdf, marginRight, startY+42, payout.Total)
//
// 	resetTextStyles(pdf)
//
// 	pdf.Line(marginLeft, startY-.5, marginRight, startY-.5)
// 	pdf.Line(marginLeft, startY+63.5, marginRight, startY+63.5)
// 	pdf.Line(297.5, startY-.5, 298.5, startY+63.5)
// }
//
// func addMonthlyPayoutTable(pdf *gopdf.GoPdf, startY float64) {
// 	setText(pdf, marginLeft, startY, "Tranzacție")
// 	setText(pdf, 328, startY, "Preț brut")
// 	setText(pdf, 424.5, startY, "Taxă Stripe")
// 	setText(pdf, 532, startY, "Total")
//
// 	pdf.Line(marginLeft, startY+21.5, marginRight, startY+21.5)
// }
//
// func addMonthlyPayoutProduct(pdf *gopdf.GoPdf, item models.Payout, startY float64) {
// 	setText(pdf, marginLeft, startY+16, item.IssueDate)
//
// 	setRightAlignedText(pdf, 367, startY, item.Gross)
// 	setRightAlignedText(pdf, 474, startY, item.StripeFees)
// 	setRightAlignedText(pdf, marginRight, startY, item.Total)
//
// 	pdf.SetTextColor(0, 0, 0)
// 	setText(pdf, marginLeft, startY, item.PayoutID)
// 	pdf.SetTextColor(94, 100, 112)
// }
//
// func (s DocumentService) GenerateDocument() error {
// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
// 	pdf.AddPage()
//
// 	err := pdf.AddTTFFont("Roboto", "./internal/pdf/static/fonts/Roboto-Regular.ttf")
// 	if err != nil {
// 		return err
// 	}
// 	err = pdf.AddTTFFont("Roboto-Bold", "./internal/pdf/static/fonts/Roboto-Bold.ttf")
// 	if err != nil {
// 		return err
// 	}
//
// 	const (
// 		itemHeight             = 50
// 		firstPageStartY        = 357.0
// 		secondPageStartY       = 135.0
// 		firstPageCapacity      = 8
// 		subsequentPageCapacity = 12
// 		firstPageTableY        = 315
// 		subsequentPageTableY   = 93
// 	)
//
// 	itemsLength := len(s.MonthlyPayout.Items)
// 	pagesNeeded := pagesNeeded(itemsLength)
// 	currentPage := 1
//
// 	resetTextStyles(&pdf)
//
// 	err = s.addHeader(&pdf, s.MonthlyPayout.IssueDate)
// 	if err != nil {
// 		return err
// 	}
// 	err = s.addFooter(&pdf, pagesNeeded, currentPage)
// 	if err != nil {
// 		return err
// 	}
// 	addMonthlyPayoutSummary(&pdf, s.MonthlyPayout)
// 	addMonthlyPayoutTable(&pdf, firstPageTableY)
//
// 	currentY := firstPageStartY
// 	maxItemsPerPage := firstPageCapacity
//
// 	var itemCounter int
// 	for _, item := range s.MonthlyPayout.Items {
// 		if itemCounter == maxItemsPerPage {
// 			pdf.AddPage()
// 			currentPage++
//
// 			err = s.addSecondaryHeader(&pdf)
// 			if err != nil {
// 				return err
// 			}
// 			err = s.addFooter(&pdf, pagesNeeded, currentPage)
// 			if err != nil {
// 				return err
// 			}
//
// 			addMonthlyPayoutTable(&pdf, subsequentPageTableY)
//
// 			currentY = secondPageStartY
// 			itemCounter = 0
// 			maxItemsPerPage = subsequentPageCapacity
// 		}
// 		addMonthlyPayoutProduct(&pdf, item, currentY)
// 		currentY += itemHeight
// 		itemCounter++
// 	}
//
// 	outputDir := "./internal/pdf/output/"
// 	pdfFile := filepath.Join(outputDir, "output.pdf")
// 	return pdf.WritePdf(pdfFile)
// }
