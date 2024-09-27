package main

import (
	"flag"
	"log"

	"github.com/diother/go-invoices/config"
	"github.com/diother/go-invoices/database"
	"github.com/diother/go-invoices/internal/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	username := flag.String("username", "", "Username for the new account")
	password := flag.String("password", "", "Password for the new account")
	role := flag.String("role", "admin", "Role for the new account (default is admin)")

	flag.Parse()
	validateFlags(*username, *password, *role)

	_, _, dsn, err := config.LoadEnv()
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

	hashedPassword, err := hashPassword(*password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user := models.NewUser(*username, hashedPassword, *role)
	if err := insertUser(db, user); err != nil {
		log.Fatalf("Insert user failed: %v", err)
	}

	log.Printf("User: %v has been successfully created", *username)
}

func validateFlags(username, password, role string) {
	if username == "" || password == "" {
		log.Fatal("Both username and password must be provided")
	}
	if role != "admin" {
		log.Fatal("Role is invalid. Only admin currently allowed")
	}
}

func insertUser(db *sqlx.DB, user *models.User) error {
	query := `
    INSERT INTO users (username, password, role)
	VALUES (:username, :password, :role)
    `
	_, err := db.NamedExec(query, user)
	return err
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
