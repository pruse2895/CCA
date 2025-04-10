package db

import (
	"context"
	"time"

	"cricketApp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateSession creates a new coaching session
func (m *MongoDB) CreateSession(ctx context.Context, session *models.Session) error {
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	_, err := m.sessionCollection.InsertOne(ctx, session)
	return err
}

// GetSessionByID retrieves a session by its ID
func (m *MongoDB) GetSessionByID(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	var session models.Session
	err := m.sessionCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionsByCoach retrieves all sessions for a specific coach
func (m *MongoDB) GetSessionsByCoach(ctx context.Context, coachID primitive.ObjectID) ([]*models.Session, error) {
	cursor, err := m.sessionCollection.Find(ctx, bson.M{"coachId": coachID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*models.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetAllSessions retrieves all sessions
func (m *MongoDB) GetAllSessions(ctx context.Context) ([]*models.Session, error) {
	cursor, err := m.sessionCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []*models.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// UpdateSession updates an existing session
func (m *MongoDB) UpdateSession(ctx context.Context, id primitive.ObjectID, session *models.Session) error {
	session.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":       session.Title,
			"description": session.Description,
			"date":        session.Date,
			"startTime":   session.StartTime,
			"endTime":     session.EndTime,
			"venue":       session.Venue,
			"maxStudents": session.MaxStudents,
			"updatedAt":   session.UpdatedAt,
		},
	}

	result, err := m.sessionCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// DeleteSession deletes a session by its ID
func (m *MongoDB) DeleteSession(ctx context.Context, id primitive.ObjectID) error {
	result, err := m.sessionCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// GetSessionsByCoachID retrieves all sessions for a specific coach
func (m *MongoDB) GetSessionsByCoachID(ctx context.Context, coachID primitive.ObjectID) ([]*models.Session, error) {
	var sessions []*models.Session
	cursor, err := m.sessionCollection.Find(ctx, bson.M{"coachId": coachID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}
