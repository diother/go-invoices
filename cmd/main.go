package main

import (
	"log"
	"net/http"

	"github.com/diother/go-invoices/config"
	"github.com/diother/go-invoices/database"
	"github.com/gorilla/mux"

	"github.com/diother/go-invoices/internal/documents"
	"github.com/diother/go-invoices/internal/handlers"
	"github.com/diother/go-invoices/internal/middleware"
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
	pwaRepo := repository.NewPWARepository(db)
	authRepo := repository.NewAuthRepository(db)

	donationService := services.NewDonationService(webhookRepo)
	payoutService := services.NewPayoutService(webhookRepo)
	documentService := documents.NewDocumentService()
	accountingService := services.NewAccountingService(pwaRepo, documentService)
	authService := services.NewAuthService(authRepo)

	// payouts := []*stripe.Payout{
	// 	{ID: "po_1PkFJUDXCtuWOFq8DYodF1nZ", Status: "paid"},
	// 	{ID: "po_1PZ0YuDXCtuWOFq8wiLw72fu", Status: "paid"},
	// 	{ID: "po_1Pj9wxDXCtuWOFq8lSKRH9Jx", Status: "paid"},
	// }
	// for _, payout := range payouts {
	// 	if err = payoutService.ProcessPayout(payout); err != nil {
	// 		return
	// 	}
	// }

	m := middleware.NewMiddleware(authService)

	webhookHandler := handlers.NewWebhookHandler(donationService, payoutService, stripeEndpointSecret)
	pwaHandler := handlers.NewPWAHandler(accountingService)
	authHandler := handlers.NewAuthHandler(authService)

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(m.Gzip(m.CacheControl(http.StripPrefix("/static/", fs))))

	router.Handle("/favicon.ico", m.CacheControl(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(w, r, "./static/images/favicon.ico")
	})))

	router.HandleFunc("/webhook", webhookHandler.HandleWebhooks).Methods("POST")

	router.Handle("/login", m.Gzip(http.HandlerFunc(authHandler.HandleLogin)))

	router.Handle("/", m.Gzip(m.HandleSessions(http.HandlerFunc(pwaHandler.HandleDashboard)))).Methods("GET")
	router.Handle("/document", m.Gzip(m.HandleSessions(http.HandlerFunc(pwaHandler.HandleDocuments)))).Methods("GET")
	router.Handle("/monthly", m.Gzip(m.HandleSessions(http.HandlerFunc(pwaHandler.HandleMonthly)))).Methods("GET")

	log.Println("Server listening at port 8080")
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
