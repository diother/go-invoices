package main

import (
	"log"

	"github.com/diother/go-invoices/internal/models"
	"github.com/diother/go-invoices/internal/pdf"
)

func main() {

	pay := models.Payout{
		IssueDate:  "01 Sep, 2024",
		PayoutDate: "15 Sep, 2024",
		PayoutID:   "po_1PZ0rmDXCtuWOFq8n33WSnN9",
		Gross:      "2000.00 lei",
		StripeFees: "200.00 lei",
		Total:      "1800.00 lei",
		Items: []models.PayoutTransaction{
			{ProductName: "Donation to Project A", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN1", Gross: "100.00 lei", StripeFee: "10.00 lei", Total: "90.00 lei"},
			{ProductName: "Donation to Project B", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN2", Gross: "150.00 lei", StripeFee: "15.00 lei", Total: "135.00 lei"},
			{ProductName: "Charity Auction Item", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN3", Gross: "200.00 lei", StripeFee: "20.00 lei", Total: "180.00 lei"},
			{ProductName: "Monthly Subscription", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN4", Gross: "300.00 lei", StripeFee: "30.00 lei", Total: "270.00 lei"},
		},
	}
	payout := pdf.PayoutPdf{Payout: &pay}
	err := payout.GenerateDocument()

	// mon := models.MonthlyPayout{
	// 	IssueDate:    "01 Sep, 2024",
	// 	ReportPeriod: "1 Jul, 2024 - 31 Jul, 2024",
	// 	Gross:        "5000.00 lei",
	// 	StripeFees:   "500.00 lei",
	// 	Total:        "4500.00 lei",
	// 	Items: []models.Payout{
	// 		{
	// 			IssueDate:  "10 Aug, 2024",
	// 			PayoutDate: "12 Aug, 2024",
	// 			PayoutID:   "po_1PZ0rmDXCtuWOFq8n33WSnN10",
	// 			Gross:      "1500.00 lei",
	// 			StripeFees: "150.00 lei",
	// 			Total:      "1350.00 lei",
	// 			Items: []models.PayoutTransaction{
	// 				{ProductName: "Donation A", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnA1", Gross: "500.00 lei", StripeFee: "50.00 lei", Total: "450.00 lei"},
	// 				{ProductName: "Donation B", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnA2", Gross: "1000.00 lei", StripeFee: "100.00 lei", Total: "900.00 lei"},
	// 			},
	// 		},
	// 		{
	// 			IssueDate:  "20 Aug, 2024",
	// 			PayoutDate: "22 Aug, 2024",
	// 			PayoutID:   "po_1PZ0rmDXCtuWOFq8n33WSnN11",
	// 			Gross:      "2000.00 lei",
	// 			StripeFees: "200.00 lei",
	// 			Total:      "1800.00 lei",
	// 			Items: []models.PayoutTransaction{
	// 				{ProductName: "Charity Auction", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnB1", Gross: "1500.00 lei", StripeFee: "150.00 lei", Total: "1350.00 lei"},
	// 				{ProductName: "Merchandise Sale", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnB2", Gross: "500.00 lei", StripeFee: "50.00 lei", Total: "450.00 lei"},
	// 			},
	// 		},
	// 		{
	// 			IssueDate:  "31 Aug, 2024",
	// 			PayoutDate: "02 Sep, 2024",
	// 			PayoutID:   "po_1PZ0rmDXCtuWOFq8n33WSnN12",
	// 			Gross:      "1500.00 lei",
	// 			StripeFees: "150.00 lei",
	// 			Total:      "1350.00 lei",
	// 			Items: []models.PayoutTransaction{
	// 				{ProductName: "Event Sponsorship", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnC1", Gross: "1000.00 lei", StripeFee: "100.00 lei", Total: "900.00 lei"},
	// 				{ProductName: "Workshop Fee", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnC2", Gross: "500.00 lei", StripeFee: "50.00 lei", Total: "450.00 lei"},
	// 			},
	// 		},
	// 	},
	// }
	// monthlyPayout := pdf.MonthlyPayoutPdf{MonthlyPayout: &mon}
	// err := monthlyPayout.GenerateDocument()

	// inv := models.Invoice{
	// 	ClientName:    "Ungureanu Daniel",
	// 	IssueDate:     "12 Aug, 2024",
	// 	TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
	// 	ProductName:   "Donație unică de 10 lei",
	// 	UnitPrice:     "10.00 lei",
	// 	Total:         "10.00 lei",
	// }
	// invoice := pdf.InvoicePdf{Invoice: &inv}
	// err := invoice.GenerateDocument()

	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
