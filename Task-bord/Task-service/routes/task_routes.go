package routes

import (
	"net/http"
	"time"

	"Task-service/db"
	"Task-service/models"
	"Task-service/rabbitmq"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTask handles POST /tasks
func CreateTask(c *gin.Context) {
	var task models.Task

	// 1. Parse request body
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Set task metadata
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = "OPEN"

	// 3. Connect to DB
	// 3. Connect to DB
	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed", "details": err.Error()})
		return
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database("taskboard").Collection("tasks")

	// 4. Insert task into MongoDB
	_, err = collection.InsertOne(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Convert ObjectIDs to strings
	var assignedTo []string
	for _, id := range task.AssignedTo {
		assignedTo = append(assignedTo, id.Hex())
	}

	// 5. Publish event to RabbitMQ (EVENT-DRIVEN COMMUNICATION)
	rabbitmq.PublishTaskCreated(
		task.ID.Hex(),
		task.Title,
		assignedTo,
	)

	// 6. Return response
	c.JSON(http.StatusCreated, task)
}

// GetTasks handles GET /tasks
func GetTasks(c *gin.Context) {
	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database("taskboard").Collection("tasks")

	// Build filter from query params
	filter := bson.M{}
	status := c.Query("status")
	if status != "" {
		filter["status"] = status
	}
	priority := c.Query("priority")
	if priority != "" {
		filter["priority"] = priority
	}
	assignedTo := c.Query("assignedTo")
	if assignedTo != "" {
		// assignedTo is stored as a list of ObjectIDs.
		// We need to match if the array contains this ID.
		objID, err := primitive.ObjectIDFromHex(assignedTo)
		if err == nil {
			filter["assignedTo"] = objID
		}
	}

	// Find options for sorting
	findOptions := options.Find()
	sortBy := c.Query("sortBy")
	if sortBy == "priority" {
		// Custom sort for priority might be tricky without integer values (High>Med>Low).
		// For strings "high", "medium", "low", alphabetical might not work.
		// For now simple sort, or better: Client can sort.
		// Let's assume sorting by CreatedAt descending by default.
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	} else if sortBy == "date" {
		findOptions.SetSort(bson.D{{Key: "due_date", Value: 1}})
	} else {
		// Default sort by CreatedAt desc
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer cursor.Close(ctx)

	tasks := []models.Task{} // Initialize as empty slice to return [] instead of null
	for cursor.Next(ctx) {
		var task models.Task
		cursor.Decode(&task)
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTaskStatus handles PUT /tasks/:id/status
func UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var body struct {
		Status    string `json:"status"`
		UpdatedBy string `json:"updatedBy"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database("taskboard").Collection("tasks")

	update := bson.M{
		"$set": bson.M{
			"status":    body.Status,
			"updatedAt": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Retrieve task title for notification (optional optimization: fetch before update)
	// For simplicity, just sending ID and Status or fetch again.
	// Let's fetch the task to get the title.
	var updatedTask models.Task
	collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedTask)

	// Convert ObjectIDs to strings
	var assignedTo []string
	for _, id := range updatedTask.AssignedTo {
		assignedTo = append(assignedTo, id.Hex())
	}

	rabbitmq.PublishTaskStatusUpdated(taskID, updatedTask.Title, body.Status, assignedTo, body.UpdatedBy, time.Now())

	c.JSON(http.StatusOK, gin.H{"message": "Task status updated"})
}

// DeleteTask handles DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	client, ctx, cancel, err := db.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer cancel()
	defer client.Disconnect(ctx)

	collection := client.Database("taskboard").Collection("tasks")

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
