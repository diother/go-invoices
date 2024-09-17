package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

type DonationService interface{}
type PayoutService interface{}

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
		log.Println("Request not of POST type")
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

	if event.Type != "checkout.session.completed" {
		log.Println("Unsupported event type")
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}

	var checkoutSession stripe.CheckoutSession
	err = json.Unmarshal(event.Data.Raw, &checkoutSession)
	if err != nil {
		log.Printf("Invalid JSON: %v\n", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
