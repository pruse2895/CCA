package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"cricketApp/models"
)

// --- Cricketer operations ---

func (m *MongoDB) GetCricketerByID(ctx context.Context, id primitive.ObjectID) (*models.Cricketer, error) {
	var cricketer models.Cricketer
	err := m.getCollection("cricketers").FindOne(ctx, bson.M{"_id": id}).Decode(&cricketer)
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

	result, err := m.getCollection("cricketers").UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // Return specific error if not found
	}
	return nil
}

func (m *MongoDB) CreateCricketer(ctx context.Context, cricketer *models.Cricketer) error {
	_, err := m.getCollection("cricketers").InsertOne(ctx, cricketer)
	// Consider checking for duplicate key errors here if needed, although indexes should handle it.
	return err
}

func (m *MongoDB) GetCricketerByEmail(ctx context.Context, email string) (*models.Cricketer, error) {
	var cricketer models.Cricketer
	err := m.getCollection("cricketers").FindOne(ctx, bson.M{"email": email}).Decode(&cricketer)
	if err != nil {
		return nil, err
	}
	return &cricketer, nil
}

func (m *MongoDB) GetCricketerByMobile(ctx context.Context, mobile string) (*models.Cricketer, error) {
	var cricketer models.Cricketer
	err := m.getCollection("cricketers").FindOne(ctx, bson.M{"mobile": mobile}).Decode(&cricketer)
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

	cursor, err := m.getCollection("cricketers").Find(ctx, bson.M{}, findOptions)
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
