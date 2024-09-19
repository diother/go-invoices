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
	donation DonationService
	payout   PayoutService
}

func NewWebhookHandler(donation DonationService, payout PayoutService) *WebhookHandler {
	return &WebhookHandler{
		donation: donation,
		payout:   payout,
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

	endpointSecret := "whsec_f49cf6c9295d48e6f73148435e254dec38d407fa911013a86b070725bc74469d"
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("Invalid signature: %v\n", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "charge.updated":
		var charge stripe.Charge
		err = json.Unmarshal(event.Data.Raw, &charge)
		if err != nil {
			log.Printf("Invalid JSON: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err = h.donation.ProcessDonation(&charge)
		if err != nil {
			log.Printf("Problem with the servers: %v\n", err)
			http.Error(w, "Problem with the servers", http.StatusInternalServerError)
			return
		}

	case "payout.reconciliation_completed":
		var payout stripe.Payout
		err = json.Unmarshal(event.Data.Raw, &payout)
		if err != nil {
			log.Printf("Invalid JSON: %v\n", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		err = h.payout.ProcessPayout(&payout)
		if err != nil {
			log.Printf("Problem with the servers: %v\n", err)
			http.Error(w, "Problem with the servers", http.StatusInternalServerError)
			return
		}

	default:
		log.Println("Unsupported event type:", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
