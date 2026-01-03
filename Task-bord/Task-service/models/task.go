package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task defines the structure of a task
type Task struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string               `json:"title" bson:"title"`
	Description string               `json:"description" bson:"description"`
	AssignedTo  []primitive.ObjectID `json:"assignedTo" bson:"assignedTo"`
	Status      string               `json:"status" bson:"status"`
	Priority    string               `json:"priority" bson:"priority"` // "low", "medium", "high"
	DueDate     time.Time            `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}
