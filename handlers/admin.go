package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"cricketApp/middleware/authmiddleware"
)

func (h *CricketerHandler) HandleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Find admin by email using the interface
	admin, err := h.db.GetAdminByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	claims := map[string]interface{}{
		"sub":  admin.ID,
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	_, tokenString, err := authmiddleware.TokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"admin": map[string]interface{}{
			"id":    admin.ID,
			"email": admin.Email,
			"name":  admin.Name,
		},
	})
}
