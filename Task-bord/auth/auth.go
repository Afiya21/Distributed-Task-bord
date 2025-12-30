package auth

import (
	"fmt"
	"regexp"
	"task-board/db"
	"task-board/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Hash the password before saving it
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Validate the email format using regex
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// Register a new user
func RegisterUser(email, password, role string) (*models.User, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("users")
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// Convert the ObjectID to a string
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return &user, nil
}
