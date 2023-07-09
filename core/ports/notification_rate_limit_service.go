package ports

//go:generate mockgen -destination=mocks/mock_notification_rate_limit_service.go -package=ports -source=notification_rate_limit_service.go NotificationRateLimitService
type NotificationRateLimitService interface {
	// IsExceeded check the rate limit configured by the given userID and notificationType.
	// If no configuration found for notificationType it will returns false
	// If the quantity of notifications sended are lower than the limit configured then will return false otherwrise it will return true
	IsExceeded(userID string, notificationType string) bool
}
