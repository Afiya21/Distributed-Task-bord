package auth

import (
	"fmt"
	"regexp"
	"task-board/db"
	"task-board/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// Compare password with the hashed password stored in DB
func comparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Generate JWT Token
func generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 1 day expiry
	})

	secretKey := []byte("your-secret-key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Login a user by checking email and password
func LoginUser(email, password string) (string, error) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("users")
	var user models.User

	// Find the user by email
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no user found with this email")
		}
		return "", err
	}

	// Compare the provided password with the stored hashed password
	if !comparePassword(user.Password, password) {
		return "", fmt.Errorf("invalid password")
	}

	// Generate JWT token after successful login
	token, err := generateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
