package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cricketApp/models"
)

// --- Announcement operations ---

// CreateAnnouncement inserts an announcement and returns the created document with its ID.
func (m *MongoDB) CreateAnnouncement(ctx context.Context, announcement *models.Announcement) (*models.Announcement, error) {
	announcement.CreatedAt = time.Now() // Set timestamp before inserting
	result, err := m.announcementCollection.InsertOne(ctx, announcement)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		// This should ideally not happen if MongoDB is working correctly
		return nil, mongo.ErrNilDocument // Or a custom error indicating ID issue
	}

	// Fetch the newly created document to get all fields including the generated ID
	var createdAnnouncement models.Announcement
	err = m.announcementCollection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&createdAnnouncement)
	if err != nil {
		// Log the error, but maybe return the original input data? Or return the error.
		// Depending on desired behavior if fetch fails after successful insert.
		return nil, err
	}

	return &createdAnnouncement, nil
}

func (m *MongoDB) GetAllAnnouncements(ctx context.Context) ([]models.Announcement, error) {
	var announcements []models.Announcement
	// Sort by createdAt descending
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := m.announcementCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &announcements); err != nil {
		return nil, err
	}

	if announcements == nil {
		return []models.Announcement{}, nil // Return empty slice, not nil
	}

	return announcements, nil
}
