package config

import (
	"fmt"
	"os"
)

func LoadEnv() (string, string, string, error) {
	stripeKey := os.Getenv("STRIPE_SECRET")
	stripeEndpointSecret := os.Getenv("STRIPE_ENDPOINT_SECRET")
	dsn := os.Getenv("DSN")

	if stripeKey == "" {
		return "", "", "", fmt.Errorf("Stripe secret is missing")
	}
	if stripeEndpointSecret == "" {
		return "", "", "", fmt.Errorf("Stripe endpoint secret is missing")
	}
	if dsn == "" {
		return "", "", "", fmt.Errorf("Database DSN is missing")
	}

	return stripeKey, stripeEndpointSecret, dsn, nil
}
