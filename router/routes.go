package router

import (
	"log"
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
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Accept", "Authorization", "if-modified-since", "Access-Control-Allow-Origin", "id"}, // Allow all headers
		// ExposedHeaders:   []string{"Link"},
		// AllowCredentials: true,
		// MaxAge:           300,
		// Debug:            true,
	}))

	// Add error handling middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panic recovered: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	// Create coach handler
	coachHandler := handlers.NewCoachHandler(database)

	// Create session handler
	sessionHandler := handlers.NewSessionHandler(database)

	// Create registration handler
	registrationHandler := handlers.NewRegistrationHandler(database)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/signup", cricketerHandler.HandleCricketerSignup) // done
		r.Post("/api/login", cricketerHandler.HandleCricketerLogin)   //done
		r.Post("/api/admin/login", cricketerHandler.HandleAdminLogin) //done
		r.Post("/api/coach/login", coachHandler.HandleCoachLogin)     //done
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.Authenticator) // JWT authentication

		// Cricketer routes
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.ValidateCricketer(database))
			r.Route("/api/cricketer", func(r chi.Router) {
				r.Get("/profile", cricketerHandler.GetCricketerProfile)    //done
				r.Put("/profile", cricketerHandler.UpdateCricketerProfile) //done
				r.Get("/announcement", cricketerHandler.GetAnnouncements)  //done
			})
		})

		// Coach routes
		r.Route("/api/coach", func(r chi.Router) {
			r.Use(authmiddleware.Authorizer("coach"))
			r.Route("/coach", func(r chi.Router) {
				r.Get("/profile", coachHandler.GetCoachProfile) //done
				r.Get("/{coachId}/sessions", sessionHandler.GetSessionsByCoach)
			})
		})

		// Admin routes
		r.Route("/api/admin", func(r chi.Router) {
			r.Use(authmiddleware.Authorizer("admin"))

			// announcements routes
			r.Route("/announcements", func(r chi.Router) {
				r.Post("/", cricketerHandler.CreateAnnouncement) //done
			})

			// coaches routes
			r.Route("/coaches", func(r chi.Router) {
				r.Post("/", coachHandler.CreateCoach)  //done
				r.Get("/", coachHandler.GetAllCoaches) //done
				r.Put("/", coachHandler.UpdateCoach)   //done
			})

			// Cricketer routes
			r.Route("/cricketers", func(r chi.Router) {
				r.Get("/", cricketerHandler.GetAllCricketers)
				r.Put("/{id}/joining-date", cricketerHandler.UpdateCricketerJoiningDate)
				r.Put("/{id}/inactive-status", cricketerHandler.UpdateCricketerInactiveStatus)
				r.Put("/{id}/due-date", cricketerHandler.UpdateCricketerDueDate)
			})

			// Session routes
			r.Route("/sessions", func(r chi.Router) {
				r.Get("/", sessionHandler.GetAllSessions)
				r.Get("/coach/{coachId}", sessionHandler.GetSessionsByCoachID)
				r.Post("/", sessionHandler.CreateSession)
				r.Get("/{id}", sessionHandler.GetSession)
				r.Put("/{id}", sessionHandler.UpdateSession)
				r.Delete("/{id}", sessionHandler.DeleteSession)
			})

			// Registration routes
			r.Route("/registrations", func(r chi.Router) {
				r.Post("/", registrationHandler.CreateRegistration)
				r.Get("/", registrationHandler.GetAllRegistrations)
				r.Get("/{id}", registrationHandler.GetRegistration)
				r.Put("/{id}", registrationHandler.UpdateRegistration)

			})

		})

	})

	return r
}
