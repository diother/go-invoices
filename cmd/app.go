package main

import (
	"fmt"
	"log"
	"os"

	"github.com/diother/go-invoices/database"
	"github.com/diother/go-invoices/internal/models"
	"github.com/diother/go-invoices/internal/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/balancetransaction"
)

func main() {
	dsn := os.Getenv("DSN")
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// addCharge(newCharge, db)
	// addPayout(newPayout, db)
	getStripe()
	fmt.Println("addition succeeded")

}

func addCharge(c models.Charge, db *sqlx.DB) {
	chargeRepo := repository.NewChargeRepository(db)

	err := chargeRepo.Insert(c)
	if err != nil {
		log.Fatalf("Error inserting: %v", err)
	}
}

func addPayout(p models.Payout, db *sqlx.DB) {
	payoutRepo := repository.NewPayoutRepository(db)

	err := payoutRepo.Insert(p)
	if err != nil {
		log.Fatalf("Error inserting: %v", err)
	}
}

func getStripe() {
	stripe.Key = "sk_test_51PVfUvDXCtuWOFq8ADmnd1iQEONLKIC6p1m1tALD67I6Ew4gRgOjoYGR7B5XK8hN0uc7iLE2Mbl9BedtgLIQubXU00XWzh1hmB"

	const balanceLimit = 5

	var balances []*stripe.BalanceTransaction
	params := &stripe.BalanceTransactionListParams{}
	// params.Limit = stripe.Int64(balanceLimit) // objects per request (max 100)
	// params.Payout = stripe.String("po_1PZ0Y9DXCtuWOFq8xapMbzlX")
	params.Type = stripe.String("payout")
	params.AddExpand("data.source")

	res := balancetransaction.List(params)

	for i := 0; i < balanceLimit && res.Next(); i++ {
		balances = append(balances, res.BalanceTransaction())
	}

	for _, balance := range balances {
		fmt.Println("----")
		fmt.Println("Type:", balance.Type)
		fmt.Println("ID:", balance.ID)
		fmt.Println("Created:", balance.Created)
		fmt.Println("Amount:", balance.Amount)
		fmt.Println("Fee:", balance.Fee)
		fmt.Println("Net:", balance.Net)
		// fmt.Println("ClientName:", balance.Source.Charge.BillingDetails.Name)
		// fmt.Println("ClientEmail:", balance.Source.Charge.BillingDetails.Email)
	}

	fmt.Println(len(balances))
}
