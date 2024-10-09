package middleware

import (
	"context"
	"net/http"

	"github.com/diother/go-invoices/internal/models"
)

type AuthService interface {
	ValidateSession(sessionToken string) (*models.User, error)
}

type Middleware struct {
	service AuthService
}

func NewMiddleware(service AuthService) *Middleware {
	return &Middleware{service: service}
}

func (m *Middleware) HandleSessions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		sessionToken := cookie.Value
		user, err := m.service.ValidateSession(sessionToken)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
