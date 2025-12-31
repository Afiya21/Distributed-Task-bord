package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Message   string             `json:"message" bson:"message"`
	IsRead    bool               `json:"is_read" bson:"is_read"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
