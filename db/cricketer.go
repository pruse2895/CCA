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

// --- Cricketer operations ---

func (m *MongoDB) GetCricketerByID(ctx context.Context, id primitive.ObjectID) (*models.Cricketer, error) {
	var cricketer models.Cricketer
	err := m.cricketerCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&cricketer)
	if err != nil {
		return nil, err
	}
	return &cricketer, nil
}

// UpdateCricketer updates specific fields of a cricketer.
func (m *MongoDB) UpdateCricketer(ctx context.Context, id primitive.ObjectID, name *string, email *string, mobile *string, hashedPassword *string) error {
	updateFields := bson.M{}
	if name != nil {
		updateFields["name"] = *name
	}
	if email != nil {
		// Consider adding uniqueness checks here if necessary before updating
		updateFields["email"] = *email
	}
	if mobile != nil {
		// Consider adding uniqueness checks here if necessary before updating
		updateFields["mobile"] = *mobile
	}
	if hashedPassword != nil {
		updateFields["password"] = *hashedPassword
	}

	if len(updateFields) == 0 {
		return nil // No fields to update
	}

	update := bson.M{"$set": updateFields}

	result, err := m.cricketerCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // Return specific error if not found
	}
	return nil
}

func (m *MongoDB) CreateCricketer(ctx context.Context, cricketer *models.Cricketer) error {
	_, err := m.cricketerCollection.InsertOne(ctx, cricketer)
	// Consider checking for duplicate key errors here if needed, although indexes should handle it.
	return err
}

func (m *MongoDB) GetCricketerByEmail(ctx context.Context, email string) (*models.Cricketer, error) {
	var cricketer models.Cricketer
	err := m.cricketerCollection.FindOne(ctx, bson.M{"email": email}).Decode(&cricketer)
	if err != nil {
		return nil, err
	}
	return &cricketer, nil
}

func (m *MongoDB) GetCricketerByMobile(ctx context.Context, mobile string) (*models.Cricketer, error) {
	var cricketer models.Cricketer
	err := m.cricketerCollection.FindOne(ctx, bson.M{"mobile": mobile}).Decode(&cricketer)
	if err != nil {
		return nil, err
	}
	return &cricketer, nil
}

func (m *MongoDB) GetAllCricketers(ctx context.Context) ([]models.Cricketer, error) {
	var cricketers []models.Cricketer
	// TODO: Add filtering, sorting, pagination options if needed
	findOptions := options.Find()
	// Example: findOptions.SetSort(bson.D{{Key: "name", Value: 1}}) // Sort by name ascending

	cursor, err := m.cricketerCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &cricketers); err != nil {
		return nil, err
	}

	if cricketers == nil {
		return []models.Cricketer{}, nil // Return empty slice, not nil
	}

	return cricketers, nil
}

// UpdateCricketerJoiningDate updates the joining date of a cricketer and sets the due date
func (m *MongoDB) UpdateCricketerJoiningDate(ctx context.Context, id primitive.ObjectID, joiningDate *time.Time) error {
	// Calculate due date (1 month after joining date)
	var dueDate *time.Time
	if joiningDate != nil {
		calculatedDueDate := joiningDate.AddDate(0, 1, 0) // Add 1 month
		dueDate = &calculatedDueDate
	}

	update := bson.M{"$set": bson.M{
		"joiningDate": joiningDate,
		"dueDate":     dueDate,
	}}
	result, err := m.cricketerCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// UpdateCricketerDueDate updates the due date of a cricketer
func (m *MongoDB) UpdateCricketerDueDate(ctx context.Context, id primitive.ObjectID, dueDate *time.Time) error {
	update := bson.M{"$set": bson.M{"dueDate": dueDate}}
	result, err := m.cricketerCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// UpdateCricketerInactiveStatus updates the inactive cricketer status
func (m *MongoDB) UpdateCricketerInactiveStatus(ctx context.Context, id primitive.ObjectID, isInactive bool) error {
	update := bson.M{"$set": bson.M{"inactiveCricketer": isInactive}}
	result, err := m.cricketerCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
