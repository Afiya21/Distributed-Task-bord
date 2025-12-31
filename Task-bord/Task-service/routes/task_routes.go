package routes

import (
	"Task-service/db"
	"Task-service/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTask handles the creation of a new task
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.CreatedAt = time.Now().Format(time.RFC3339)
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("tasks")

	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Task created", "task_id": result.InsertedID})
}

// GetTasks retrieves all tasks
func GetTasks(c *gin.Context) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()

	collection := client.Database("taskboard").Collection("tasks")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tasks)
}

// UpdateTask updates an existing task
// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	taskID := c.Param("id") // Get the task ID from URL parameter
	var task models.Task    // Define the task model to hold the updated data

	// Bind the incoming JSON data to the task struct
	if err := c.BindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	client, ctx, cancel := db.ConnectDB() // Get database connection
	defer cancel()                        // Ensure the cancel function is called after DB operations

	collection := client.Database("taskboard").Collection("tasks")

	// Convert taskID string to MongoDB ObjectID type
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID format"})
		return
	}

	filter := bson.M{"_id": objectID} // MongoDB filter to find task by ID
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"assigned_to": task.AssignedTo,
			"status":      task.Status,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
		},
	}

	// Perform the update operation
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// If no documents were updated, inform the user
	if updateResult.MatchedCount == 0 {
		c.JSON(404, gin.H{"message": "Task not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Task updated successfully"})
}

// DeleteTask deletes a task by ID
// DeleteTask deletes a task by ID
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id") // Get the task ID from URL parameter

	client, ctx, cancel := db.ConnectDB() // Get database connection
	defer cancel()                        // Ensure the cancel function is called after DB operations

	collection := client.Database("taskboard").Collection("tasks")

	// Convert taskID string to MongoDB ObjectID type
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid task ID format"})
		return
	}

	// Perform the delete operation
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}
