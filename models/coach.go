package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coach struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Mobile    string             `json:"mobile" bson:"mobile" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required,min=6"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
}

type UpdateCoachRequest struct {
	Name     string `json:"name,omitempty"`
	IsActive *bool  `json:"isActive,omitempty"`
}
