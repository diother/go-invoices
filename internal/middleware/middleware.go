package middleware

import (
	"compress/gzip"
	"context"
	"net/http"
	"strings"

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

func (m *Middleware) CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzip.NewWriter(w)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		gw := gzipResponseWriter{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(gw, r)

		if err := gw.Writer.Close(); err != nil {
			http.Error(w, "Failed to compress response", http.StatusInternalServerError)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
