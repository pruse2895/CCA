package models

type Admin struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`
}
