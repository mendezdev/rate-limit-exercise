package domain

import "time"

type Notification struct {
	UserID           string    `json:"user_id"`
	NotificationType string    `json:"notification_type"`
	DateCreated      time.Time `json:"date_created"`
	Message          string    `json:"message"`
}

type NotificationRequest struct {
	UserID           string `json:"user_id"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
}
