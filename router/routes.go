package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"cricketApp/db"
	"cricketApp/handlers"
	"cricketApp/middleware/authmiddleware"
)

func SetupRouter(database db.Database, cricketerHandler *handlers.CricketerHandler) http.Handler {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)         // Request logging
	r.Use(middleware.Recoverer)      // Panic recovery
	r.Use(middleware.RealIP)         // Get real IP
	r.Use(middleware.RequestID)      // Add request ID
	r.Use(middleware.Heartbeat("/")) // Health check endpoint

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/signup", cricketerHandler.HandleCricketerSignup)
		r.Post("/api/login", cricketerHandler.HandleCricketerLogin)
		r.Post("/api/admin/login", cricketerHandler.HandleAdminLogin)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.Authenticator) // JWT authentication

		// Cricketer routes
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.ValidateCricketer(database)) // Pass database to middleware factory
			r.Route("/api/cricketer", func(r chi.Router) {
				r.Get("/profile", cricketerHandler.GetCricketerProfile)
				r.Put("/profile", cricketerHandler.UpdateCricketerProfile)
			})
		})

		// Admin routes
		r.Route("/api/admin", func(r chi.Router) {
			r.Use(authmiddleware.Authorizer("admin")) // Only admin can access
			r.Get("/cricketers", cricketerHandler.GetAllCricketers)
			r.Post("/announcements", cricketerHandler.CreateAnnouncement)
		})

		// Announcement routes (accessible to all authenticated users)
		r.Route("/api/announcements", func(r chi.Router) {
			r.Get("/", cricketerHandler.GetAnnouncements)
		})
	})

	return r
}
