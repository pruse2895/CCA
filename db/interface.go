package db

import (
	"context"
	"cricketApp/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Database interface defines all database operations
type Database interface {
	// Cricketer operations
	GetCricketerByID(ctx context.Context, id primitive.ObjectID) (*models.Cricketer, error)
	UpdateCricketer(ctx context.Context, id primitive.ObjectID, name *string, email *string, mobile *string, hashedPassword *string) error
	CreateCricketer(ctx context.Context, cricketer *models.Cricketer) error
	GetCricketerByEmail(ctx context.Context, email string) (*models.Cricketer, error)
	GetCricketerByMobile(ctx context.Context, mobile string) (*models.Cricketer, error)
	GetAllCricketers(ctx context.Context) ([]models.Cricketer, error)
	UpdateCricketerJoiningDate(ctx context.Context, id primitive.ObjectID, joiningDate *time.Time) error
	UpdateCricketerDueDate(ctx context.Context, id primitive.ObjectID, dueDate *time.Time) error
	UpdateCricketerInactiveStatus(ctx context.Context, id primitive.ObjectID, isInactive bool) error

	// Coach operations
	CreateCoach(ctx context.Context, coach *models.Coach) error
	GetCoachByMobile(ctx context.Context, mobile string) (*models.Coach, error)
	GetCoachByID(ctx context.Context, id primitive.ObjectID) (*models.Coach, error)
	GetAllCoaches(ctx context.Context) ([]models.Coach, error)
	UpdateCoach(ctx context.Context, id primitive.ObjectID, coach *models.Coach) error

	// Admin operations
	GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error)
	GetAdminByID(ctx context.Context, id primitive.ObjectID) (*models.Admin, error)

	// Announcement operations
	CreateAnnouncement(ctx context.Context, announcement *models.Announcement) (*models.Announcement, error)
	GetAllAnnouncements(ctx context.Context) ([]models.Announcement, error)

	// Session methods
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByID(ctx context.Context, id primitive.ObjectID) (*models.Session, error)
	GetSessionsByCoach(ctx context.Context, coachID primitive.ObjectID) ([]*models.Session, error)
	GetAllSessions(ctx context.Context) ([]*models.Session, error)
	UpdateSession(ctx context.Context, id primitive.ObjectID, session *models.Session) error
	DeleteSession(ctx context.Context, id primitive.ObjectID) error
	GetSessionsByCoachID(ctx context.Context, coachID primitive.ObjectID) ([]*models.Session, error)

	// Registration methods
	CreateRegistration(ctx context.Context, registration *models.RegistrationForm) error
	GetRegistrationByID(ctx context.Context, id primitive.ObjectID) (*models.RegistrationForm, error)
	GetAllRegistrations(ctx context.Context) ([]*models.RegistrationForm, error)
	UpdateRegistration(ctx context.Context, id primitive.ObjectID, registration *models.RegistrationForm) error
}
