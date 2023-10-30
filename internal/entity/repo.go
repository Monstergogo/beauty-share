package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ShareInfo struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Text      string             `json:"text" bson:"text"`
	Images    []string           `json:"images" bson:"images"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
