package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session represents a coaching session
type Session struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CoachID     primitive.ObjectID `json:"coachId" bson:"coachId" binding:"required"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description"`
	Date        time.Time          `json:"date" bson:"date" binding:"required"`
	StartTime   time.Time          `json:"startTime" bson:"startTime" binding:"required"`
	EndTime     time.Time          `json:"endTime" bson:"endTime" binding:"required"`
	Venue       string             `json:"venue" bson:"venue" binding:"required"`
	MaxStudents int                `json:"maxStudents" bson:"maxStudents" binding:"required,min=1"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CreateSessionRequest represents the request body for creating a new session
type CreateSessionRequest struct {
	CoachID     string    `json:"coachId" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	StartTime   time.Time `json:"startTime" binding:"required"`
	EndTime     time.Time `json:"endTime" binding:"required"`
	Venue       string    `json:"venue" binding:"required"`
	MaxStudents int       `json:"maxStudents" binding:"required,min=1"`
}

// UpdateSessionRequest represents the request body for updating a session
type UpdateSessionRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	StartTime   *time.Time `json:"startTime,omitempty"`
	EndTime     *time.Time `json:"endTime,omitempty"`
	Venue       *string    `json:"venue,omitempty"`
	MaxStudents *int       `json:"maxStudents,omitempty"`
}
