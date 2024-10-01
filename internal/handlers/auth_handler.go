package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/diother/go-invoices/internal/custom_errors"
	"github.com/diother/go-invoices/internal/helpers"
	"github.com/diother/go-invoices/internal/models"
)

type AuthService interface {
	Authenticate(user, password string) (*models.User, error)
	GenerateSession(user *models.User) (*models.Session, error)
}

type AuthHandler struct {
	service AuthService
	tmpl    *template.Template
}

func NewAuthHandler(service AuthService) *AuthHandler {
	tmpl := template.New("base").Funcs(template.FuncMap{
		"arr": helpers.ComponentHelper,
	})
	tmpl, err := tmpl.ParseGlob("internal/views/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	tmpl, err = tmpl.ParseGlob("internal/views/components/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	return &AuthHandler{
		service: service,
		tmpl:    tmpl,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := h.tmpl.ExecuteTemplate(w, "login", nil); err != nil {
			log.Printf("Template execution failed: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %v", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := h.service.Authenticate(username, password)
		if err != nil {
			var credentialsError *custom_errors.CredentialsError
			if errors.As(err, &credentialsError) {
				http.Error(w, credentialsError.Error(), http.StatusUnauthorized)
				return
			}
			log.Printf("Auth service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusBadRequest)
			return
		}

		session, err := h.service.GenerateSession(user)
		if err != nil {
			log.Printf("Auth service error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusBadRequest)
		}

		cookie := http.Cookie{
			Name:     "session_token",
			Value:    fmt.Sprintf("%d", session.SessionToken),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Unix(session.ExpiresAt, 0),
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	return
}
