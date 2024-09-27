package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/diother/go-invoices/internal/dto"
	"github.com/signintech/gopdf"
)

type AccountingService interface {
	GenerateInvoice(id string) (*gopdf.GoPdf, error)
	GeneratePayoutReport(id string) (*gopdf.GoPdf, error)
	GenerateMonthlyReport(date string) (*gopdf.GoPdf, error)
	GenerateMonthlyReportView(date string) (*dto.MonthlyReportView, error)
}

type PWAHandler struct {
	service AccountingService
}

func NewPWAHandler(service AccountingService) *PWAHandler {
	return &PWAHandler{service: service}
}

func (h *PWAHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("home").ParseGlob("internal/views/*.html"))
	if err := tmpl.ExecuteTemplate(w, "home", nil); err != nil {
		log.Printf("Template execution failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *PWAHandler) HandleDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	documentType := r.FormValue("type")
	documentID := r.FormValue("ID")
	documentDate := r.FormValue("date")

	if err := validateDocumentRequest(documentType, documentID, documentDate); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var pdf *gopdf.GoPdf
	var err error

	switch documentType {
	case "donation":
		pdf, err = h.service.GenerateInvoice(documentID)
		if err != nil {
			log.Printf("Accounting service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusBadRequest)
			return
		}

	case "payout":
		pdf, err = h.service.GeneratePayoutReport(documentID)
		if err != nil {
			log.Printf("Accounting service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusBadRequest)
			return
		}

	case "monthly":
		pdf, err = h.service.GenerateMonthlyReport(documentDate)
		if err != nil {
			log.Printf("Accounting service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusBadRequest)
			return
		}

	default:
		http.Error(w, "Invalid document type", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=document.pdf")

	if _, err = pdf.WriteTo(w); err != nil {
		http.Error(w, "Failed to write PDF", http.StatusInternalServerError)
		return
	}
}

func (h *PWAHandler) HandleMonthly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	documentDate := r.FormValue("date")
	if documentDate == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	data, err := h.service.GenerateMonthlyReportView(documentDate)
	if err != nil {
		log.Printf("Accounting service error: %v\n", err)
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	var buffer bytes.Buffer
	tmpl := template.Must(template.New("monthly").ParseGlob("internal/views/*.html"))
	if err := tmpl.ExecuteTemplate(&buffer, "monthly", data); err != nil {
		log.Printf("Template execution failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	buffer.WriteTo(w)
}

func validateDocumentRequest(documentType, documentID, documentDate string) error {
	if documentType == "" {
		return fmt.Errorf("")
	}
	switch documentType {
	case "monthly":
		if documentDate == "" {
			return fmt.Errorf("")
		}
	default:
		if documentID == "" {
			return fmt.Errorf("")
		}
	}
	return nil
}
