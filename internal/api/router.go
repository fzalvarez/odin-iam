package api

import (
	"net/http"

	"github.com/fzalvarez/odin-iam/internal/api/handlers"
	"github.com/fzalvarez/odin-iam/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterParams struct {
	AuthService *auth.AuthService
}

func NewRouter(p RouterParams) http.Handler {
	r := chi.NewRouter()

	// Middlewares básicos
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth handler
	authHandler := handlers.NewAuthHandler(p.AuthService)

	// Auth endpoints
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)

	return r
}
