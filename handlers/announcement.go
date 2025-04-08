package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"

	"cricketApp/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateAnnouncement is now a method of CricketerHandler
func (h *CricketerHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announcement models.Announcement
	if err := json.NewDecoder(r.Body).Decode(&announcement); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get admin ID (subject) from JWT claims
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	adminID, ok := claims["sub"].(string) // Assuming admin ID (string) is in sub
	if !ok {
		http.Error(w, "Invalid token claims (sub is not string)", http.StatusUnauthorized)
		return
	}

	// Validate if the adminID exists in the database
	_, err = h.db.GetAdminByID(r.Context(), adminID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Admin user not found", http.StatusUnauthorized) // Or Forbidden
		} else {
			http.Error(w, "Error validating admin user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	announcement.CreatedBy = adminID

	// Use the database interface, now returns the created object
	createdAnnouncement, err := h.db.CreateAnnouncement(r.Context(), &announcement)
	if err != nil {
		http.Error(w, "Error creating announcement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Return the announcement object received from the db layer (includes ID)
	json.NewEncoder(w).Encode(createdAnnouncement)
}

// GetAnnouncements is now a method of CricketerHandler
func (h *CricketerHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	// Use the database interface
	announcements, err := h.db.GetAllAnnouncements(r.Context())
	if err != nil {
		http.Error(w, "Error fetching announcements: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(announcements)
}
