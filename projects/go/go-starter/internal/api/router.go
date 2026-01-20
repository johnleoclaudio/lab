package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/go-starter/internal/api/handlers"
	"github.com/yourusername/go-starter/internal/api/middleware"
	"github.com/yourusername/go-starter/internal/db"
	"github.com/yourusername/go-starter/internal/repository"
	"github.com/yourusername/go-starter/internal/service"
	"log/slog"
)

func NewRouter(queries *db.Queries, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recovery(logger))
	// CORS middleware can be added here if needed
	// r.Use(middleware.CORS(allowedOrigins, allowedMethods, allowedHeaders))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Initialize dependencies (following clean architecture)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, logger)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", userHandler.GetUser)
		})
	})

	return r
}
