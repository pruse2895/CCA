package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cricketApp/models"
)

// --- Admin operations ---

func (m *MongoDB) CreateAdmin(ctx context.Context, admin *models.Admin) error {
	_, err := m.adminCollection.InsertOne(ctx, admin)
	return err
}

func (m *MongoDB) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	err := m.adminCollection.FindOne(ctx, bson.M{"email": email}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (m *MongoDB) GetAdminByID(ctx context.Context, id primitive.ObjectID) (*models.Admin, error) {
	var admin models.Admin
	err := m.adminCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (m *MongoDB) UpdateAdmin(ctx context.Context, id primitive.ObjectID, admin *models.Admin) error {
	update := bson.M{"$set": admin}
	result, err := m.adminCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
