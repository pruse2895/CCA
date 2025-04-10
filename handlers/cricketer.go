package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"cricketApp/db"
	"cricketApp/middleware/authmiddleware"
	"cricketApp/models"
)

// CricketerHandler holds the database interface
type CricketerHandler struct {
	db db.Database
}

// NewCricketerHandler creates a new CricketerHandler
func NewCricketerHandler(db db.Database) *CricketerHandler {
	return &CricketerHandler{db: db}
}

func (h *CricketerHandler) HandleCricketerSignup(w http.ResponseWriter, r *http.Request) {
	var cricketer models.Cricketer
	if err := json.NewDecoder(r.Body).Decode(&cricketer); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cricketer.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	cricketer.Password = string(hashedPassword)
	cricketer.ID = primitive.NewObjectID() // Assign a new ObjectID
	cricketer.CreatedAt = time.Now()       // Set creation timestamp

	// Check if email or mobile already exists using the interface
	_, err = h.db.GetCricketerByEmail(r.Context(), cricketer.Email)
	if err == nil {
		// Found existing email
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		// Other database error
		http.Error(w, "Database error checking email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for mobile number
	_, err = h.db.GetCricketerByMobile(r.Context(), cricketer.Mobile)
	if err == nil {
		// Found existing mobile
		http.Error(w, "Mobile number already exists", http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		// Other database error
		http.Error(w, "Database error checking mobile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert new cricketer using the interface
	err = h.db.CreateCricketer(r.Context(), &cricketer)
	if err != nil {
		http.Error(w, "Error creating cricketer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cricketer created successfully"})
}

func (h *CricketerHandler) HandleCricketerLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find cricketer by mobile number using the interface
	cricketer, err := h.db.GetCricketerByMobile(r.Context(), loginRequest.Mobile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid mobile number or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error fetching cricketer: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(cricketer.Password), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Invalid mobile number or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	claims := map[string]interface{}{
		"sub":  cricketer.ID.Hex(),
		"role": "cricketer",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	_, tokenString, err := authmiddleware.TokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"cricketer": map[string]interface{}{
			"id":     cricketer.ID.Hex(),
			"name":   cricketer.Name,
			"email":  cricketer.Email,
			"mobile": cricketer.Mobile,
		},
	})
}

func (h *CricketerHandler) GetCricketerProfile(w http.ResponseWriter, r *http.Request) {
	// Get cricketer ID from JWT claims (assuming 'sub' holds the ObjectID hex)
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
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

	// Fetch fresh data from database using the interface
	cricketer, err := h.db.GetCricketerByID(r.Context(), cricketerID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Cricketer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching profile: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return profile without sensitive information
	profile := map[string]interface{}{
		"id":                cricketer.ID.Hex(),
		"name":              cricketer.Name,
		"email":             cricketer.Email,
		"mobile":            cricketer.Mobile,
		"createdAt":         cricketer.CreatedAt,
		"joiningDate":       cricketer.JoiningDate,
		"dueDate":           cricketer.DueDate,
		"inactiveCricketer": cricketer.InactiveCricketer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *CricketerHandler) UpdateCricketerProfile(w http.ResponseWriter, r *http.Request) {
	// Get cricketer ID from JWT claims
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
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

	// Decode update request
	var updateData struct {
		Name     *string `json:"name,omitempty"` // Use pointers to handle omitted fields
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare fields for update (handle password hashing)
	var hashedPassPtr *string
	if updateData.Password != nil && *updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		hashedPassStr := string(hashedPassword)
		hashedPassPtr = &hashedPassStr
	}

	// Update in database using the interface
	err = h.db.UpdateCricketer(r.Context(), cricketerID,
		updateData.Name, updateData.Email, nil, hashedPassPtr,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Cricketer not found", http.StatusNotFound)
		} else {
			// Check for potential duplicate key errors if the interface/implementation signals them
			// e.g., if err == db.ErrDuplicateEmail { ... }
			http.Error(w, "Error updating profile: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Fetch updated profile to return (optional, could just return success message)
	updatedCricketer, err := h.db.GetCricketerByID(r.Context(), cricketerID)
	if err != nil {
		log.Printf("Warning: Failed to fetch updated profile after update for %s: %v", cricketerIDHex, err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully, but could not fetch updated details."})
		return
	}

	// Return updated profile without sensitive information
	profile := map[string]interface{}{
		"id":                updatedCricketer.ID.Hex(),
		"name":              updatedCricketer.Name,
		"email":             updatedCricketer.Email,
		"mobile":            updatedCricketer.Mobile,
		"createdAt":         updatedCricketer.CreatedAt,
		"joiningDate":       updatedCricketer.JoiningDate,
		"dueDate":           updatedCricketer.DueDate,
		"inactiveCricketer": updatedCricketer.InactiveCricketer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// GetAllCricketers fetches all cricketer profiles (admin only)
func (h *CricketerHandler) GetAllCricketers(w http.ResponseWriter, r *http.Request) {
	cricketers, err := h.db.GetAllCricketers(r.Context())
	if err != nil {
		http.Error(w, "Error fetching cricketers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// IMPORTANT: Map to a response model to avoid exposing sensitive data like passwords
	// Define a response struct or use map[string]interface{}
	responseProfiles := make([]map[string]interface{}, len(cricketers))
	for i, c := range cricketers {
		responseProfiles[i] = map[string]interface{}{
			"id":                c.ID.Hex(),
			"name":              c.Name,
			"email":             c.Email,
			"mobile":            c.Mobile,
			"createdAt":         c.CreatedAt,
			"joiningDate":       c.JoiningDate,
			"dueDate":           c.DueDate,
			"inactiveCricketer": c.InactiveCricketer,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseProfiles)
}

// UpdateCricketerJoiningDate updates the joining date of a cricketer (admin only)
func (h *CricketerHandler) UpdateCricketerJoiningDate(w http.ResponseWriter, r *http.Request) {
	// Parse cricketer ID from URL
	cricketerIDHex := chi.URLParam(r, "id")
	fmt.Println("Cricketer ID:", cricketerIDHex)
	cricketerID, err := primitive.ObjectIDFromHex(cricketerIDHex)
	if err != nil {
		http.Error(w, "Invalid cricketer ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		JoiningDate time.Time `json:"joiningDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the joining date in the database
	err = h.db.UpdateCricketerJoiningDate(r.Context(), cricketerID, &request.JoiningDate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Cricketer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating joining date: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Joining date updated successfully"})
}

// UpdateCricketerInactiveStatus updates whether a cricketer is inactive or not (admin only)
func (h *CricketerHandler) UpdateCricketerInactiveStatus(w http.ResponseWriter, r *http.Request) {
	// Parse cricketer ID from URL
	fmt.Println("Path:", r.URL.Path)
	fmt.Println("Param ID:", chi.URLParam(r, "id"))
	cricketerIDHex := chi.URLParam(r, "id")
	cricketerID, err := primitive.ObjectIDFromHex(cricketerIDHex)
	if err != nil {
		http.Error(w, "Invalid cricketer ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request struct {
		IsInactive bool `json:"isInactive"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the inactive status in the database
	err = h.db.UpdateCricketerInactiveStatus(r.Context(), cricketerID, request.IsInactive)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Cricketer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating inactive status: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Inactive status updated successfully"})
}
