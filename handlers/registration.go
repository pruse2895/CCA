package handlers

import (
	"encoding/json"
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
	var req models.CreateRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create new registration form
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
	}

	// Convert CricketerID string to ObjectID
	cricketerID, err := primitive.ObjectIDFromHex(req.CricketerID)
	if err != nil {
		http.Error(w, "Invalid cricketer ID format", http.StatusBadRequest)
		return
	}
	registration.CricketerID = cricketerID

	if err := h.db.CreateRegistration(r.Context(), registration); err != nil {
		http.Error(w, "Failed to create registration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Registration created successfully",
		"registration": registration,
	})
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
	if updateData.Status != nil {
		registration.Status = *updateData.Status
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
