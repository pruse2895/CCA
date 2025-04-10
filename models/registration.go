package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParentDetails represents the parent/guardian information
type ParentDetails struct {
	Name       string `json:"name" bson:"name" binding:"required"`
	ContactNo  string `json:"contactNo" bson:"contactNo" binding:"required"`
	Occupation string `json:"occupation" bson:"occupation" binding:"required"`
}

// RegistrationForm represents the complete registration form
type RegistrationForm struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FormNo           string             `json:"formNo" bson:"formNo" binding:"required"`
	Date             string             `json:"date" bson:"date" binding:"required"`
	Reference        string             `json:"reference" bson:"reference"`
	FullName         string             `json:"fullName" bson:"fullName" binding:"required"`
	DateOfBirth      string             `json:"dateOfBirth" bson:"dateOfBirth" binding:"required"`
	ResidenceAddress string             `json:"residenceAddress" bson:"residenceAddress" binding:"required"`
	ContactNo        string             `json:"contactNo" bson:"contactNo" binding:"required"`
	Email            string             `json:"email" bson:"email" binding:"required,email"`
	Education        string             `json:"education" bson:"education" binding:"required"`
	SchoolCollege    string             `json:"schoolCollege" bson:"schoolCollege" binding:"required"`
	AadhaarNo        string             `json:"aadhaarNo" bson:"aadhaarNo" binding:"required"`
	Whatsapp         string             `json:"whatsapp" bson:"whatsapp" binding:"required"`
	ParentDetails    ParentDetails      `json:"parentDetails" bson:"parentDetails" binding:"required"`
	CricketerID      primitive.ObjectID `json:"cricketerId,omitempty" bson:"cricketerId,omitempty"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CreateRegistrationRequest represents the request body for creating a new registration
type CreateRegistrationRequest struct {
	FormNo           string        `json:"formNo" binding:"required"`
	Date             string        `json:"date" binding:"required"`
	Reference        string        `json:"reference"`
	FullName         string        `json:"fullName" binding:"required"`
	DateOfBirth      string        `json:"dateOfBirth" binding:"required"`
	ResidenceAddress string        `json:"residenceAddress" binding:"required"`
	ContactNo        string        `json:"contactNo" binding:"required"`
	Email            string        `json:"email" binding:"required,email"`
	Education        string        `json:"education" binding:"required"`
	SchoolCollege    string        `json:"schoolCollege" binding:"required"`
	AadhaarNo        string        `json:"aadhaarNo" binding:"required"`
	Whatsapp         string        `json:"whatsapp" binding:"required"`
	ParentDetails    ParentDetails `json:"parentDetails" binding:"required"`
	CricketerID      string        `json:"cricketerId" binding:"required"`
}

// UpdateRegistrationRequest represents the request body for updating a registration
type UpdateRegistrationRequest struct {
	FormNo           *string        `json:"formNo"`
	Date             *string        `json:"date"`
	Reference        *string        `json:"reference"`
	FullName         *string        `json:"fullName"`
	DateOfBirth      *string        `json:"dateOfBirth"`
	ResidenceAddress *string        `json:"residenceAddress"`
	ContactNo        *string        `json:"contactNo"`
	Email            *string        `json:"email"`
	Education        *string        `json:"education"`
	SchoolCollege    *string        `json:"schoolCollege"`
	AadhaarNo        *string        `json:"aadhaarNo"`
	Whatsapp         *string        `json:"whatsapp"`
	ParentDetails    *ParentDetails `json:"parentDetails"`
}
