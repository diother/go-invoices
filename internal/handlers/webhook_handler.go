package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

type DonationService interface {
	ProcessDonation(charge *stripe.Charge) error
}

type PayoutService interface {
	ProcessPayout(payout *stripe.Payout) error
}

type WebhookHandler struct {
	donation       DonationService
	payout         PayoutService
	endpointSecret string
}

func NewWebhookHandler(donation DonationService, payout PayoutService, secret string) *WebhookHandler {
	return &WebhookHandler{
		donation:       donation,
		payout:         payout,
		endpointSecret: secret,
	}
}

func (h *WebhookHandler) HandleWebhooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.endpointSecret)
	if err != nil {
		log.Printf("Invalid signature: %v\n", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "charge.updated":
		var charge stripe.Charge
		if err = json.Unmarshal(event.Data.Raw, &charge); err != nil {
			log.Printf("Invalid JSON: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err = h.donation.ProcessDonation(&charge); err != nil {
			log.Printf("Service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	case "payout.reconciliation_completed":
		var payout stripe.Payout
		if err = json.Unmarshal(event.Data.Raw, &payout); err != nil {
			log.Printf("Invalid JSON: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err = h.payout.ProcessPayout(&payout); err != nil {
			log.Printf("Service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	default:
		log.Println("Unsupported event type:", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
