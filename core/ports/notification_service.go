package ports

import "github.com/mendezdev/rate-limit-example/core/domain"

//go:generate mockgen -destination=mocks/mock_notification_service.go -package=ports -source=notification_service.go NotificationService

type NotificationService interface {
	Send(domain.NotificationRequest) error
}
