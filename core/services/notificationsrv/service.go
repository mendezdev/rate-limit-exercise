package notificationsrv

import (
	"fmt"
	"time"

	"github.com/mendezdev/rate-limit-example/core/domain"
	"github.com/mendezdev/rate-limit-example/core/ports"
	"github.com/mendezdev/rate-limit-example/core/utils"
)

type service struct {
	mailSender                 ports.MailSenderGateway
	notificationRepo           ports.NotificationRepository
	rateLimitConfigurationRepo ports.RateLimitConfigurationRepository
	timeProvider               utils.TimeProvider
}

func New(
	mailSender ports.MailSenderGateway,
	notificationRepo ports.NotificationRepository,
	rateLimitConfigurationRepo ports.RateLimitConfigurationRepository,
	timeProvider utils.TimeProvider,
) ports.NotificationService {
	return service{
		mailSender:                 mailSender,
		notificationRepo:           notificationRepo,
		rateLimitConfigurationRepo: rateLimitConfigurationRepo,
		timeProvider:               timeProvider,
	}
}

// Send implements ports.NotificationService.
func (srv service) Send(nr domain.NotificationRequest) error {
	if srv.isRateLimitExceededFor(nr.UserID, nr.NotificationType) {
		return fmt.Errorf("the rate limit has been exceeded for user_id %s and notification type %s", nr.UserID, nr.NotificationType)
	}

	if err := srv.mailSender.Send(nr.UserID, nr.Message); err != nil {
		return fmt.Errorf("error trying to send email: %s", err.Error())
	}

	// create and save Notification
	notification := domain.Notification{
		NotificationType: nr.NotificationType,
		UserID:           nr.UserID,
		DateCreated:      srv.timeProvider.Now(),
		Message:          nr.Message,
	}

	if err := srv.notificationRepo.Save(notification); err != nil {
		fmt.Printf("error trying to save notification: %s", err.Error())
	}

	return nil
}

// isRateLimitExceededFor check the rate limit configured by the given userID and notificationType.
// If no configuration found for notificationType it will returns false
// If the quantity of notifications sended are lower than the limit configured then will return false otherwrise it will return true
func (srv service) isRateLimitExceededFor(userID string, notificationType string) bool {
	rlc, err := srv.rateLimitConfigurationRepo.GetRateLimitConfiguration(notificationType)
	if err != nil {
		fmt.Printf("error trying to get rate limit configuration: %s\n", err.Error())
		return false
	}

	if rlc == nil {
		fmt.Printf("no rate limit configuration was found for notificaiton type: %s\n", notificationType)
		return false
	}

	from := srv.timeProvider.
		Now().
		Add(-time.Duration(rlc.TimeUnit) * rlc.GetTimeMeasureInDuration())

	notifications, err := srv.notificationRepo.GetByTypeAndUserAndFromDate(userID, notificationType, from)

	if err != nil {
		fmt.Printf("error trying to get notifications: %s\n", err.Error())
		return false
	}

	return len(notifications) >= rlc.Limit
}
