package rabbitmq

import (
	"common/rabbitmq"
	"log"
	"time"
)

var (
	RabbitClient *rabbitmq.RabbitClient
)

func InitRabbitMQ(url string) {
	var err error
	RabbitClient, err = rabbitmq.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
}

func PublishTaskCreated(taskId, title string, assignedTo []string) {
	if RabbitClient == nil {
		log.Println("RabbitMQ client not initialized")
		return
	}

	payload := map[string]interface{}{
		"taskId":     taskId,
		"title":      title,
		"assignedTo": assignedTo,
	}

	// Publish TaskCreated (for Audit/Logs) - Ignoring separate audit queue for now, or just log?
	// Let's send to task_notifications too just in case consumer wants it, or just drop it if no queue.
	// Actually user-service/notification-service don't listen to task_events.
	// Let's comment out task.created for now or send to task_notifications?
	// The plan said notification service handles TaskAssigned.
	// Let's just fix TaskAssigned.

	// Publish TaskAssigned (for Notification)
	RabbitClient.Publish("", "task_notifications", "TaskAssigned", payload)
}

func PublishTaskStatusUpdated(taskId, title, status string, assignedTo []string, updatedBy string, updatedAt time.Time) {
	if RabbitClient == nil {
		return
	}
	payload := map[string]interface{}{
		"taskId":     taskId,
		"title":      title,
		"status":     status,
		"assignedTo": assignedTo,
		"updatedBy":  updatedBy,
		"updatedAt":  updatedAt,
	}
	RabbitClient.Publish("", "task_notifications", "TaskStatusUpdated", payload)
}
