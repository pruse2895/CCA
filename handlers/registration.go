package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cricketApp/db"
	"cricketApp/models"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegistrationHandler struct {
	db db.Database
}

func NewRegistrationHandler(db db.Database) *RegistrationHandler {
	return &RegistrationHandler{db: db}
}

// CreateRegistration handles the creation of a new registration form
func (h *RegistrationHandler) CreateRegistration(w http.ResponseWriter, r *http.Request) {
	// Read and log the raw request body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Debug - Handler: Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Debug - Handler: Raw request body: %s\n", string(body))

	// Reset the body reader so it can be read again
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.CreateRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Debug - Handler: Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.FormNo == "" || req.FullName == "" || req.ContactNo == "" || req.Date == "" || req.DateOfBirth == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Convert cricketerId from string to ObjectID
	cricketerID, err := primitive.ObjectIDFromHex(req.CricketerID)
	if err != nil {
		http.Error(w, "Invalid cricketer ID", http.StatusBadRequest)
		return
	}

	// Check if registration already exists for this cricketer
	existingRegistrations, err := h.db.GetAllRegistrations(r.Context())
	if err != nil {
		http.Error(w, "Error checking existing registrations", http.StatusInternalServerError)
		return
	}

	for _, reg := range existingRegistrations {
		if reg.CricketerID == cricketerID {
			http.Error(w, "A registration already exists for this cricketer", http.StatusConflict)
			return
		}
	}

	// Create registration form
	registration := &models.RegistrationForm{
		FormNo:           req.FormNo,
		Date:             req.Date,
		Reference:        req.Reference,
		FullName:         req.FullName,
		DateOfBirth:      req.DateOfBirth,
		ResidenceAddress: req.ResidenceAddress,
		ContactNo:        req.ContactNo,
		Email:            req.Email,
		Education:        req.Education,
		SchoolCollege:    req.SchoolCollege,
		AadhaarNo:        req.AadhaarNo,
		Whatsapp:         req.Whatsapp,
		ParentDetails:    req.ParentDetails,
		CricketerID:      cricketerID,
	}

	// Save to database
	err = h.db.CreateRegistration(r.Context(), registration)
	if err != nil {
		http.Error(w, "Error creating registration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration created successfully"})
}

// GetRegistration retrieves a registration by ID
func (h *RegistrationHandler) GetRegistration(w http.ResponseWriter, r *http.Request) {
	registrationID := chi.URLParam(r, "id")
	if registrationID == "" {
		http.Error(w, "Registration ID is required", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(registrationID)
	if err != nil {
		http.Error(w, "Invalid registration ID", http.StatusBadRequest)
		return
	}

	registration, err := h.db.GetRegistrationByID(r.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Registration not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching registration", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registration)
}

// GetAllRegistrations retrieves all registrations
func (h *RegistrationHandler) GetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	registrations, err := h.db.GetAllRegistrations(r.Context())
	if err != nil {
		http.Error(w, "Error fetching registrations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registrations)
}

// UpdateRegistration updates an existing registration
func (h *RegistrationHandler) UpdateRegistration(w http.ResponseWriter, r *http.Request) {
	registrationID := chi.URLParam(r, "id")
	if registrationID == "" {
		http.Error(w, "Registration ID is required", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(registrationID)
	if err != nil {
		http.Error(w, "Invalid registration ID", http.StatusBadRequest)
		return
	}

	// Get existing registration
	registration, err := h.db.GetRegistrationByID(r.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Registration not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching registration", http.StatusInternalServerError)
		}
		return
	}

	// Decode update request
	var updateData models.UpdateRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if updateData.FormNo != nil {
		registration.FormNo = *updateData.FormNo
	}
	if updateData.Date != nil {
		registration.Date = *updateData.Date
	}
	if updateData.Reference != nil {
		registration.Reference = *updateData.Reference
	}
	if updateData.FullName != nil {
		registration.FullName = *updateData.FullName
	}
	if updateData.DateOfBirth != nil {
		registration.DateOfBirth = *updateData.DateOfBirth
	}
	if updateData.ResidenceAddress != nil {
		registration.ResidenceAddress = *updateData.ResidenceAddress
	}
	if updateData.ContactNo != nil {
		registration.ContactNo = *updateData.ContactNo
	}
	if updateData.Email != nil {
		registration.Email = *updateData.Email
	}
	if updateData.Education != nil {
		registration.Education = *updateData.Education
	}
	if updateData.SchoolCollege != nil {
		registration.SchoolCollege = *updateData.SchoolCollege
	}
	if updateData.AadhaarNo != nil {
		registration.AadhaarNo = *updateData.AadhaarNo
	}
	if updateData.Whatsapp != nil {
		registration.Whatsapp = *updateData.Whatsapp
	}
	if updateData.ParentDetails != nil {
		registration.ParentDetails = *updateData.ParentDetails
	}

	// Update registration in database
	if err := h.db.UpdateRegistration(r.Context(), objID, registration); err != nil {
		http.Error(w, "Error updating registration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Registration updated successfully",
		"registration": registration,
	})
}
