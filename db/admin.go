package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cricketApp/models"
)

// --- Admin operations ---

func (m *MongoDB) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	err := m.getCollection("admins").FindOne(ctx, bson.M{"email": email}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (m *MongoDB) GetAdminByID(ctx context.Context, id string) (*models.Admin, error) {
	var admin models.Admin
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // Invalid ID format
	}
	err = m.getCollection("admins").FindOne(ctx, bson.M{"_id": objID}).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
