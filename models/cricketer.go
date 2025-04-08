package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cricketer struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" binding:"required"`
	Mobile   string             `json:"mobile" bson:"mobile" binding:"required"`
	Email    string             `json:"email" bson:"email" binding:"required,email"`
	Password string             `json:"password" bson:"password" binding:"required,min=6"`
}
