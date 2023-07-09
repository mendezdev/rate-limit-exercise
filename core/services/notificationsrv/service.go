package notificationsrv

import (
	"fmt"

	"github.com/mendezdev/rate-limit-example/core/domain"
	"github.com/mendezdev/rate-limit-example/core/ports"
	"github.com/mendezdev/rate-limit-example/core/utils"
)

type service struct {
	mailSender               ports.MailSenderGateway
	notificationRateLimitSrv ports.NotificationRateLimitService
	notificationRepo         ports.NotificationRepository
	timeProvider             utils.TimeProvider
}

func New(
	mailSender ports.MailSenderGateway,
	notificationRateLimitSrv ports.NotificationRateLimitService,
	notificationRepo ports.NotificationRepository,
	timeProvider utils.TimeProvider,
) ports.NotificationService {
	return service{
		mailSender:               mailSender,
		notificationRateLimitSrv: notificationRateLimitSrv,
		notificationRepo:         notificationRepo,
		timeProvider:             timeProvider,
	}
}

// Send implements ports.NotificationService.
func (srv service) Send(nr domain.NotificationRequest) error {
	if srv.notificationRateLimitSrv.IsExceeded(nr.UserID, nr.NotificationType) {
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
