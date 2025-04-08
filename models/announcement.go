package models

import (
	"time"
)

type Announcement struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedBy string    `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
