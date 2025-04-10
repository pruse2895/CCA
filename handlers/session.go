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

type SessionHandler struct {
	db db.Database
}

func NewSessionHandler(db db.Database) *SessionHandler {
	return &SessionHandler{db: db}
}

// CreateSession creates a new coaching session
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert coach ID string to ObjectID
	coachID, err := primitive.ObjectIDFromHex(req.CoachID)
	if err != nil {
		http.Error(w, "Invalid coach ID", http.StatusBadRequest)
		return
	}

	// Create new session
	session := &models.Session{
		CoachID:     coachID,
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Venue:       req.Venue,
		MaxStudents: req.MaxStudents,
	}

	if err := h.db.CreateSession(r.Context(), session); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Session created successfully",
		"session": session,
	})
}

// GetSession retrieves a session by ID
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	session, err := h.db.GetSessionByID(r.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Session not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching session", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// GetSessionsByCoach retrieves all sessions for a specific coach
func (h *SessionHandler) GetSessionsByCoach(w http.ResponseWriter, r *http.Request) {
	coachID := chi.URLParam(r, "coachId")
	if coachID == "" {
		http.Error(w, "Coach ID is required", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(coachID)
	if err != nil {
		http.Error(w, "Invalid coach ID", http.StatusBadRequest)
		return
	}

	sessions, err := h.db.GetSessionsByCoach(r.Context(), objID)
	if err != nil {
		http.Error(w, "Error fetching sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// GetAllSessions retrieves all sessions
func (h *SessionHandler) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.db.GetAllSessions(r.Context())
	if err != nil {
		http.Error(w, "Error fetching sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// UpdateSession updates an existing session
func (h *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Get existing session
	session, err := h.db.GetSessionByID(r.Context(), objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Session not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching session", http.StatusInternalServerError)
		}
		return
	}

	// Decode update request
	var updateData models.UpdateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if updateData.Title != nil {
		session.Title = *updateData.Title
	}
	if updateData.Description != nil {
		session.Description = *updateData.Description
	}
	if updateData.Date != nil {
		session.Date = *updateData.Date
	}
	if updateData.StartTime != nil {
		session.StartTime = *updateData.StartTime
	}
	if updateData.EndTime != nil {
		session.EndTime = *updateData.EndTime
	}
	if updateData.Venue != nil {
		session.Venue = *updateData.Venue
	}
	if updateData.MaxStudents != nil {
		session.MaxStudents = *updateData.MaxStudents
	}

	// Update session in database
	if err := h.db.UpdateSession(r.Context(), objID, session); err != nil {
		http.Error(w, "Error updating session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Session updated successfully",
		"session": session,
	})
}

// DeleteSession deletes a session
func (h *SessionHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteSession(r.Context(), objID); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Session not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting session", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session deleted successfully"})
}

// GetSessionsByCoachID handles getting all sessions for a specific coach
func (h *SessionHandler) GetSessionsByCoachID(w http.ResponseWriter, r *http.Request) {
	// Get coach ID from URL parameter
	coachIDHex := chi.URLParam(r, "coachId")
	if coachIDHex == "" {
		http.Error(w, "Coach ID is required", http.StatusBadRequest)
		return
	}

	// Convert coach ID from string to ObjectID
	coachID, err := primitive.ObjectIDFromHex(coachIDHex)
	if err != nil {
		http.Error(w, "Invalid coach ID", http.StatusBadRequest)
		return
	}

	// Get sessions from database
	sessions, err := h.db.GetSessionsByCoachID(r.Context(), coachID)
	if err != nil {
		http.Error(w, "Error fetching sessions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return sessions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}
