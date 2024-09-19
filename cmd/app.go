package main

import (
	"log"
	// "net/http"
	"os"

	"github.com/diother/go-invoices/database"
	"github.com/stripe/stripe-go/v79"
	// "github.com/diother/go-invoices/internal/handlers"
	"github.com/diother/go-invoices/internal/repository"
	"github.com/diother/go-invoices/internal/services"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("DSN")
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	donationRepo := repository.NewDonationRepositoryMySQL(db)
	payoutRepo := repository.NewPayoutRepositoryMySQL(db)

	// donationService := services.NewDonationServiceImpl(donationRepo)
	payoutService := services.NewPayoutServiceImpl(payoutRepo, donationRepo)

	p := &stripe.Payout{
		Status: "paid",
	}
	err = payoutService.ProcessPayout(p)
	if err != nil {
		log.Println(err)
	}

	// webhookHandler := handlers.NewWebhookHandler(donationService, payoutService)
	//
	// http.HandleFunc("/webhook", webhookHandler.HandleWebhooks)

	// log.Println("Server listening at port 8080")
	// err = http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
}
