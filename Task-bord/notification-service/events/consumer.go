package events

import (
	"common/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notification-service/db"
	"notification-service/models"
	"notification-service/websockets"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetupConsumer starts listening for task events
func SetupConsumer(url string, hub *websockets.Hub) {
	client, err := rabbitmq.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Consume for Task Events
	err = client.Consume("task_notifications", func(event rabbitmq.Event) {
		log.Printf("[Notification Service] Received event: %s", event.Type)
		processEvent(event, hub)
	})
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
}

func processEvent(event rabbitmq.Event, hub *websockets.Hub) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	collection := client.Database("taskboard").Collection("notifications")

	// Convert payload to map
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		// Try to marshal/unmarshal again if it came as raw string or map
		jsonBody, _ := json.Marshal(event.Payload)
		if err := json.Unmarshal(jsonBody, &payload); err != nil {
			log.Printf("Error processing payload: %v", err)
			return
		}
	}

	switch event.Type {
	case "TaskAssigned":
		// Payload: taskId, title, assignedTo (list of IDs)
		assignedToInterface, ok := payload["assignedTo"].([]interface{})
		if !ok {
			log.Printf("Error: assignedTo is not a list in payload: %v", payload)
			return
		}

		taskTitle := fmtToString(payload["title"])

		for _, userIdInterface := range assignedToInterface {
			userId := fmtToString(userIdInterface)
			message := "You have been assigned to task: " + taskTitle
			notification := models.Notification{
				ID:        primitive.NewObjectID(),
				UserID:    userId,
				Message:   message,
				IsRead:    false,
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			_, err := collection.InsertOne(ctx, notification) // Use ctx from ConnectDB
			if err != nil {
				log.Printf("Failed to save notification for user %s: %v", userId, err)
			} else {
				log.Printf("Notification saved for user %s", userId)
				// REAL-TIME PUSH
				hub.Broadcast <- websockets.Message{UserID: userId, Content: notification}
			}
		}

	case "TaskStatusUpdated":
		// Payload: taskId, status, assignedTo (list), title
		assignedToInterface, _ := payload["assignedTo"].([]interface{})
		status := fmtToString(payload["status"])
		taskTitle := fmtToString(payload["title"])

		updatedBy, _ := payload["updatedBy"].(string)
		updatedAtStr, _ := payload["updatedAt"].(string)

		// Notify Assigned Users
		assignedToInterface, ok := payload["assignedTo"].([]interface{})
		if ok {
			for _, userID := range assignedToInterface {
				if idStr, ok := userID.(string); ok {
					msg := fmt.Sprintf("Task '%s' status updated to %s", taskTitle, status)

					notification := models.Notification{
						ID:        primitive.NewObjectID(),
						UserID:    idStr,
						Message:   msg,
						IsRead:    false,
						CreatedAt: time.Now().Format(time.RFC3339),
					}
					collection.InsertOne(ctx, notification)
					SendNotification(hub, idStr, notification)
				}
			}
		}

		// Notify Admins
		// 1. Fetch all users from User Service (simple http call)
		resp, err := http.Get("http://localhost:8087/users")
		if err == nil {
			defer resp.Body.Close()
			var users []struct {
				ID       string `json:"id"`
				Role     string `json:"role"`
				Username string `json:"username"`
				Email    string `json:"email"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&users); err == nil {
				// Resolve Updater Name
				updaterName := "Unknown User"
				for _, u := range users {
					if u.ID == updatedBy {
						if u.Username != "" {
							updaterName = u.Username
						} else {
							updaterName = u.Email
						}
						break
					}
				}

				// Format Date
				var friendlyTime string
				parsedTime, err := time.Parse(time.RFC3339, updatedAtStr)
				if err == nil {
					friendlyTime = parsedTime.Format("Jan 02, 2006 at 3:04 PM")
				} else {
					friendlyTime = updatedAtStr // Fallback
				}

				for _, u := range users {
					if u.Role == "admin" {
						msg := fmt.Sprintf("ADMIN ALERT: Task '%s' updated to %s by %s at %s",
							taskTitle, status, updaterName, friendlyTime)

						notification := models.Notification{
							ID:        primitive.NewObjectID(),
							UserID:    u.ID,
							Message:   msg,
							IsRead:    false,
							CreatedAt: time.Now().Format(time.RFC3339),
						}
						collection.InsertOne(ctx, notification)
						SendNotification(hub, u.ID, notification)
					}
				}
			}
		}
		log.Printf(">>> NOTIFICATION [ADMIN]: Task '%v' status updated to '%v'", taskTitle, status)

	default:
		log.Printf("Unknown event type: %s", event.Type)
	}
}

func SendNotification(hub *websockets.Hub, userID string, content interface{}) {
	hub.Broadcast <- websockets.Message{UserID: userID, Content: content}
}

func fmtToString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
