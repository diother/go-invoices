package main

import (
	"log"
	"net/http"

	"github.com/diother/go-invoices/config"
	"github.com/diother/go-invoices/database"
	"github.com/diother/go-invoices/internal/handlers"
	"github.com/diother/go-invoices/internal/repository"
	"github.com/diother/go-invoices/internal/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	stripeKey, stripeEndpointSecret, dsn, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Environment variable is missing: %v", err)
	}

	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	stripe.Key = stripeKey

	donationRepo := repository.NewDonationRepositoryMySQL(db)
	payoutRepo := repository.NewPayoutRepositoryMySQL(db)

	donationService := services.NewDonationServiceImpl(donationRepo)
	payoutService := services.NewPayoutServiceImpl(payoutRepo, donationRepo)

	webhookHandler := handlers.NewWebhookHandler(donationService, payoutService, stripeEndpointSecret)
	http.HandleFunc("/webhook", webhookHandler.HandleWebhooks)

	log.Println("Server listening at port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
