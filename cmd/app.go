package main

import (
	"log"
	"os"

	"github.com/diother/go-invoices/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("DSN")
	database, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer database.Close()
}

// func getStripe() {
// 	stripe.Key = "sk_test_51PVfUvDXCtuWOFq8ADmnd1iQEONLKIC6p1m1tALD67I6Ew4gRgOjoYGR7B5XK8hN0uc7iLE2Mbl9BedtgLIQubXU00XWzh1hmB"
//
// 	const balanceLimit = 5
//
// 	var balances []*stripe.BalanceTransaction
// 	params := &stripe.BalanceTransactionListParams{}
// 	// params.Limit = stripe.Int64(balanceLimit) // objects per request (max 100)
// 	// params.Payout = stripe.String("po_1PZ0Y9DXCtuWOFq8xapMbzlX")
// 	params.Type = stripe.String("charge")
// 	params.AddExpand("data.source")
//
// 	res := balancetransaction.List(params)
//
// 	for i := 0; i < balanceLimit && res.Next(); i++ {
// 		balances = append(balances, res.BalanceTransaction())
// 	}
//
// 	for _, balance := range balances {
// 		fmt.Println("----")
// 		fmt.Println("Type:", balance.Type)
// 		fmt.Println("ID:", balance.ID)
// 		fmt.Println("Date:", balance.Created)
// 		fmt.Println("Amount:", balance.Amount)
// 		fmt.Println("Fee:", balance.Fee)
// 		fmt.Println("Net:", balance.Net)
// 		fmt.Println("Source:", balance.Source.Charge.BillingDetails.Name)
// 		fmt.Println("Source:", balance.Source.Charge.BillingDetails.Email)
// 	}
//
// 	fmt.Println(len(balances))
// }
