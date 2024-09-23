package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/diother/go-invoices/internal/models"
	"github.com/signintech/gopdf"
)

type AccountingService interface {
	FetchDonations() ([]*models.Donation, error)
	GenerateDocument(documentType, documentID string) (gopdf.GoPdf, error)
}

type PwaHandler struct {
	service AccountingService
}

func NewPwaHandler(service AccountingService) *PwaHandler {
	return &PwaHandler{service: service}
}

func (h *PwaHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	donations, err := h.service.FetchDonations()
	if err != nil {
		fmt.Println(fmt.Errorf("Select failed: %w", err))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.New("test").ParseGlob("internal/views/*.html"))
	if err = tmpl.ExecuteTemplate(w, "home", donations); err != nil {
		fmt.Println(fmt.Errorf("Select failed: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *PwaHandler) HandleDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	documentType := r.FormValue("type")
	documentID := r.FormValue("ID")
	if err := validateDocumentRequest(documentType, documentID); err != nil {
		http.Error(w, "Parameters are missing", http.StatusBadRequest)
		return
	}

	pdf, err := h.service.GenerateDocument(documentType, documentID)
	if err != nil {
		log.Printf("Accounting service error: %v\n", err)
		http.Error(w, "Parameters are missing", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=output.pdf")

	if _, err = pdf.WriteTo(w); err != nil {
		http.Error(w, "Failed to write PDF", http.StatusInternalServerError)
		return
	}
}

func validateDocumentRequest(documentType, documentID string) error {
	if documentType == "" {
		return fmt.Errorf("")
	}
	if documentID == "" {
		return fmt.Errorf("")
	}
	return nil
}
