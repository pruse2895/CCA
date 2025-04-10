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

// CreateCoach creates a new coach in the database
func (m *MongoDB) CreateCoach(ctx context.Context, coach *models.Coach) error {
	// Set creation timestamp
	coach.CreatedAt = time.Now()

	// Create unique index for mobile number if it doesn't exist
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "mobile", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := m.coachCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil && !isIndexAlreadyExistsError(err) {
		return err
	}

	// Insert the coach
	_, err = m.coachCollection.InsertOne(ctx, coach)
	return err
}

// GetCoachByMobile retrieves a coach by their mobile number
func (m *MongoDB) GetCoachByMobile(ctx context.Context, mobile string) (*models.Coach, error) {
	var coach models.Coach
	err := m.coachCollection.FindOne(ctx, bson.M{"mobile": mobile}).Decode(&coach)
	if err != nil {
		return nil, err
	}
	return &coach, nil
}

// GetCoachByID retrieves a coach by their ID
func (m *MongoDB) GetCoachByID(ctx context.Context, id primitive.ObjectID) (*models.Coach, error) {
	var coach models.Coach
	err := m.coachCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&coach)
	if err != nil {
		return nil, err
	}
	return &coach, nil
}

// GetAllCoaches retrieves all coaches
func (m *MongoDB) GetAllCoaches(ctx context.Context) ([]models.Coach, error) {
	var coaches []models.Coach
	cursor, err := m.coachCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &coaches); err != nil {
		return nil, err
	}

	if coaches == nil {
		return []models.Coach{}, nil
	}

	return coaches, nil
}

func (m *MongoDB) UpdateCoach(ctx context.Context, id primitive.ObjectID, coach *models.Coach) error {
	update := bson.M{"$set": coach}
	result, err := m.coachCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
