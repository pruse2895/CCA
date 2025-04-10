package db

import (
	"context"
	"time"

	"cricketApp/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateRegistration creates a new registration form
func (m *MongoDB) CreateRegistration(ctx context.Context, registration *models.RegistrationForm) error {
	registration.CreatedAt = time.Now()
	registration.UpdatedAt = time.Now()

	_, err := m.registrationCollection.InsertOne(ctx, registration)
	return err
}

// GetRegistrationByID retrieves a registration by its ID
func (m *MongoDB) GetRegistrationByID(ctx context.Context, id primitive.ObjectID) (*models.RegistrationForm, error) {
	var registration models.RegistrationForm
	err := m.registrationCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&registration)
	if err != nil {
		return nil, err
	}
	return &registration, nil
}

// GetAllRegistrations retrieves all registrations
func (m *MongoDB) GetAllRegistrations(ctx context.Context) ([]*models.RegistrationForm, error) {
	cursor, err := m.registrationCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var registrations []*models.RegistrationForm
	if err = cursor.All(ctx, &registrations); err != nil {
		return nil, err
	}
	return registrations, nil
}

// UpdateRegistration updates an existing registration
func (m *MongoDB) UpdateRegistration(ctx context.Context, id primitive.ObjectID, registration *models.RegistrationForm) error {
	registration.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"formNo":           registration.FormNo,
			"date":             registration.Date,
			"reference":        registration.Reference,
			"fullName":         registration.FullName,
			"dateOfBirth":      registration.DateOfBirth,
			"residenceAddress": registration.ResidenceAddress,
			"contactNo":        registration.ContactNo,
			"email":            registration.Email,
			"education":        registration.Education,
			"schoolCollege":    registration.SchoolCollege,
			"aadhaarNo":        registration.AadhaarNo,
			"whatsapp":         registration.Whatsapp,
			"parentDetails":    registration.ParentDetails,
			"cricketerId":      registration.CricketerID,
			"updatedAt":        registration.UpdatedAt,
		},
	}

	result, err := m.registrationCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
