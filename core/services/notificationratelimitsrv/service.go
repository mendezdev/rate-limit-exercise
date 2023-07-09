package notificationratelimitsrv

import (
	"fmt"
	"time"

	"github.com/mendezdev/rate-limit-example/core/ports"
	"github.com/mendezdev/rate-limit-example/core/utils"
)

type service struct {
	rateLimitConfigurationRepo ports.RateLimitConfigurationRepository
	notificationRepo           ports.NotificationRepository
	timeProvider               utils.TimeProvider
}

func New(
	rlcr ports.RateLimitConfigurationRepository,
	nr ports.NotificationRepository,
	tp utils.TimeProvider,
) ports.NotificationRateLimitService {
	return service{
		rateLimitConfigurationRepo: rlcr,
		notificationRepo:           nr,
		timeProvider:               tp,
	}
}

// IsExceeded implements ports.NotificationRateLimitService.
func (srv service) IsExceeded(userID string, notificationType string) bool {
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
