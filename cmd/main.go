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

	if err = database.ApplyMigrations(dsn); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	stripe.Key = stripeKey

	webhookRepo := repository.NewWebhookRepository(db)

	donationService := services.NewDonationService(webhookRepo)
	payoutService := services.NewPayoutService(webhookRepo)
	accountingService := services.NewAccountingService(webhookRepo)

	webhookHandler := handlers.NewWebhookHandler(donationService, payoutService, stripeEndpointSecret)
	pwaHandler := handlers.NewPwaHandler(accountingService)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/webhook", webhookHandler.HandleWebhooks)
	http.HandleFunc("/", pwaHandler.Test)
	http.HandleFunc("/document", pwaHandler.Document)

	log.Println("Server listening at port 8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
