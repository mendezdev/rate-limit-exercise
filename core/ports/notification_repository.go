package ports

import (
	"time"

	"github.com/mendezdev/rate-limit-example/core/domain"
)

//go:generate mockgen -destination=mocks/mock_notification_repository.go -package=ports -source=notification_repository.go NotificationRepository

type NotificationRepository interface {
	// GetByTypeAndUserAndFromDate will return notifications filtered by userID, notificationType and fromDate
	GetByTypeAndUserAndFromDate(userID string, notificationType string, fromDate time.Time) ([]domain.Notification, error)

	// Save will store the notification
	Save(domain.Notification) error
}
