package main

import (
	"log"

	"github.com/diother/go-invoices/models"
	"github.com/diother/go-invoices/views"
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
			{ProductName: "Fundraiser Gala Ticket", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN5", Gross: "250.00 lei", StripeFee: "25.00 lei", Total: "225.00 lei"},
			{ProductName: "Merchandise Sale", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN6", Gross: "100.00 lei", StripeFee: "10.00 lei", Total: "90.00 lei"},
			{ProductName: "Event Sponsorship", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN7", Gross: "400.00 lei", StripeFee: "40.00 lei", Total: "360.00 lei"},
			{ProductName: "Online Course Enrollment", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN8", Gross: "120.00 lei", StripeFee: "12.00 lei", Total: "108.00 lei"},
			{ProductName: "Workshop Registration", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN9", Gross: "80.00 lei", StripeFee: "8.00 lei", Total: "72.00 lei"},
			{ProductName: "Donation to Health Fund", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN0", Gross: "90.00 lei", StripeFee: "9.00 lei", Total: "81.00 lei"},
			{ProductName: "Holiday Appeal", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN10", Gross: "110.00 lei", StripeFee: "11.00 lei", Total: "99.00 lei"},
			{ProductName: "Community Support", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN11", Gross: "130.00 lei", StripeFee: "13.00 lei", Total: "117.00 lei"},
			{ProductName: "Animal Welfare Donation", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN12", Gross: "160.00 lei", StripeFee: "16.00 lei", Total: "144.00 lei"},
			{ProductName: "Scholarship Fund", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN13", Gross: "180.00 lei", StripeFee: "18.00 lei", Total: "162.00 lei"},
			{ProductName: "Environmental Project", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN14", Gross: "140.00 lei", StripeFee: "14.00 lei", Total: "126.00 lei"},
			{ProductName: "Cultural Event Donation", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN15", Gross: "170.00 lei", StripeFee: "17.00 lei", Total: "153.00 lei"},
			{ProductName: "Support for Local Artists", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN16", Gross: "200.00 lei", StripeFee: "20.00 lei", Total: "180.00 lei"},
			{ProductName: "Donation to Project A", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN1", Gross: "100.00 lei", StripeFee: "10.00 lei", Total: "90.00 lei"},
			{ProductName: "Donation to Project B", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN2", Gross: "150.00 lei", StripeFee: "15.00 lei", Total: "135.00 lei"},
			{ProductName: "Public Health Campaign", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN17", Gross: "190.00 lei", StripeFee: "19.00 lei", Total: "171.00 lei"},
			{ProductName: "Emergency Relief Fund", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN18", Gross: "210.00 lei", StripeFee: "21.00 lei", Total: "189.00 lei"},
			{ProductName: "Education Support", TransactionId: "ch_1PZ0rmDXCtuWOFq8n33WSnN19", Gross: "230.00 lei", StripeFee: "23.00 lei", Total: "207.00 lei"},
		},
	}
	payout := views.PayoutView{Payout: &pay}
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
	// monthlyPayout := views.MonthlyPayoutView{MonthlyPayout: &mon}
	// err := monthlyPayout.GenerateDocument()

	// inv := models.Invoice{
	// 	ClientName:    "Ungureanu Daniel",
	// 	IssueDate:     "12 Aug, 2024",
	// 	TransactionId: "pi_3Pn0hXDXCtuWOFq820psOpql",
	// 	ProductName:   "Donație unică de 10 lei",
	// 	UnitPrice:     "10.00 lei",
	// 	Total:         "10.00 lei",
	// }
	// invoice := views.InvoiceView{Invoice: &inv}
	// err := invoice.GenerateDocument()

	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
	}
}
