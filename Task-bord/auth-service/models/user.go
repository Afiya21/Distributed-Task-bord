package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
	Username string             `bson:"username"` // Add username field
}
