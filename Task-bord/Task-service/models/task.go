package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task defines the structure of a task
type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	AssignedTo  primitive.ObjectID `json:"assigned_to" bson:"assigned_to"`
	Status      string             `json:"status" bson:"status"`
	DueDate     string             `json:"due_date" bson:"due_date"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
}
