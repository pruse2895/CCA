package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"cricketApp/db"
	"cricketApp/middleware/authmiddleware"
	"cricketApp/models"

	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type CoachHandler struct {
	db db.Database
}

func NewCoachHandler(db db.Database) *CoachHandler {
	return &CoachHandler{db: db}
}

func (h *CoachHandler) HandleCoachLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	coach, err := h.db.GetCoachByMobile(r.Context(), loginRequest.Mobile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid mobile number or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(coach.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid mobile number or password", http.StatusUnauthorized)
		return
	}

	claims := map[string]interface{}{
		"sub":  coach.ID.Hex(),
		"role": "coach",
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
		"coach": map[string]interface{}{
			"id":     coach.ID.Hex(),
			"name":   coach.Name,
			"mobile": coach.Mobile,
		},
	})
}

func (h *CoachHandler) GetCoachProfile(w http.ResponseWriter, r *http.Request) {
	// Get coach ID from JWT claims
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	coachIDHex, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Invalid token subject (sub)", http.StatusUnauthorized)
		return
	}
	coachID, err := primitive.ObjectIDFromHex(coachIDHex)
	if err != nil {
		http.Error(w, "Invalid coach ID format in token", http.StatusUnauthorized)
		return
	}

	// Fetch coach profile
	coach, err := h.db.GetCoachByID(r.Context(), coachID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Coach not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching profile", http.StatusInternalServerError)
		}
		return
	}

	// Return profile without sensitive information
	profile := map[string]interface{}{
		"id":        coach.ID.Hex(),
		"name":      coach.Name,
		"mobile":    coach.Mobile,
		"createdAt": coach.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *CoachHandler) CreateCoach(w http.ResponseWriter, r *http.Request) {
	var coach models.Coach
	if err := json.NewDecoder(r.Body).Decode(&coach); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(coach.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	coach.Password = string(hashedPassword)
	coach.ID = primitive.NewObjectID()
	coach.CreatedAt = time.Now()

	// Check if mobile already exists
	_, err = h.db.GetCoachByMobile(r.Context(), coach.Mobile)
	if err == nil {
		http.Error(w, "Mobile number already exists", http.StatusConflict)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, "Database error checking mobile", http.StatusInternalServerError)
		return
	}

	// Create coach
	if err := h.db.CreateCoach(r.Context(), &coach); err != nil {
		http.Error(w, "Error creating coach", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Coach created successfully"})
}

func (h *CoachHandler) GetAllCoaches(w http.ResponseWriter, r *http.Request) {
	coaches, err := h.db.GetAllCoaches(r.Context())
	if err != nil {
		http.Error(w, "Error fetching coaches", http.StatusInternalServerError)
		return
	}

	// Map to response model to avoid exposing sensitive data
	responseCoaches := make([]map[string]interface{}, len(coaches))
	for i, c := range coaches {
		responseCoaches[i] = map[string]interface{}{
			"id":        c.ID.Hex(),
			"name":      c.Name,
			"mobile":    c.Mobile,
			"isActive":  c.IsActive,
			"createdAt": c.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseCoaches)
}

func (h *CoachHandler) UpdateCoach(w http.ResponseWriter, r *http.Request) {
	coachID := r.URL.Query().Get("id")
	if coachID == "" {
		http.Error(w, "Coach ID is required", http.StatusBadRequest)
		return
	}

	var updateData models.UpdateCoachRequest
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(coachID)
	if err != nil {
		http.Error(w, "Invalid coach ID", http.StatusBadRequest)
		return
	}

	// Get existing coach
	coach, err := h.db.GetCoachByID(r.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Coach not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching coach", http.StatusInternalServerError)
		}
		return
	}

	// Update fields if provided
	if updateData.Name != "" {
		coach.Name = updateData.Name
	}

	if updateData.IsActive != nil {
		coach.IsActive = *updateData.IsActive
	}

	// Update coach in database
	if err := h.db.UpdateCoach(r.Context(), objID, coach); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Coach not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating coach", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Coach updated successfully"})
}
