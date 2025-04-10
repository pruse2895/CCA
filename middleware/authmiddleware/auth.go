package authmiddleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cricketApp/db"
	// Need models for cricketer type
)

var (
	TokenAuth *jwtauth.JWTAuth
)

// Define a custom type for context keys
type contextKey string

// Define the specific key we will use
const cricketerKey contextKey = "cricketer"

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)
}

// GetEmailFromClaims extracts email from JWT claims
func GetEmailFromClaims(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", err
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return "", err
	}

	return email, nil
}

// Authenticator middleware validates JWT tokens
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from header
		tokenString := jwtauth.TokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "No token found", http.StatusUnauthorized)
			return
		}

		// Verify and decode token
		token, err := TokenAuth.Decode(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Token:", token)

		// Add token to context
		ctx := jwtauth.NewContext(r.Context(), token, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authorizer middleware checks if user has required role
func Authorizer(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())
			if err != nil || token == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Logger middleware logs request details
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// ValidateCricketer middleware checks if the cricketer exists in the database
// It now accepts a db.Database interface and returns the middleware handler.
func ValidateCricketer(database db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())
			if err != nil || token == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Use ID from claims (assuming 'sub' holds the ObjectID hex)
			cricketerIDHex, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid token subject (sub)", http.StatusUnauthorized)
				return
			}
			cricketerID, err := primitive.ObjectIDFromHex(cricketerIDHex)
			if err != nil {
				http.Error(w, "Invalid cricketer ID format in token", http.StatusUnauthorized)
				return
			}

			// Check if cricketer exists using the database interface
			cricketer, err := database.GetCricketerByID(r.Context(), cricketerID)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					http.Error(w, "Cricketer not found", http.StatusNotFound)
				} else {
					log.Printf("Database error validating cricketer %s: %v", cricketerIDHex, err) // Log the actual error
					http.Error(w, "Database error", http.StatusInternalServerError)
				}
				return
			}

			// Add cricketer to context using the custom key type
			ctx := context.WithValue(r.Context(), cricketerKey, *cricketer) // Use cricketerKey
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
