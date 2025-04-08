package db

import (
	"context"
	"cricketApp/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Admin operations
	GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error)
	GetAdminByID(ctx context.Context, id string) (*models.Admin, error)

	// Announcement operations
	CreateAnnouncement(ctx context.Context, announcement *models.Announcement) (*models.Announcement, error)
	GetAllAnnouncements(ctx context.Context) ([]models.Announcement, error)
}

// MongoDB implements the Database interface
type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(client *mongo.Client, dbName string) *MongoDB {
	return &MongoDB{
		client: client,
		db:     client.Database(dbName),
	}
}

// GetCollection returns a collection by name
func (m *MongoDB) getCollection(name string) *mongo.Collection {
	return m.db.Collection(name)
}
