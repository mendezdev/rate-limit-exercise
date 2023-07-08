package application

import (
	"fmt"

	"github.com/mendezdev/rate-limit-example/core/domain"
	"github.com/mendezdev/rate-limit-example/core/ports"
	"github.com/mendezdev/rate-limit-example/core/services/notificationsrv"
	"github.com/mendezdev/rate-limit-example/core/utils"
	"github.com/mendezdev/rate-limit-example/insfrastructure/gateway"
	"github.com/mendezdev/rate-limit-example/insfrastructure/repositories/notificationrepo"
	"github.com/mendezdev/rate-limit-example/insfrastructure/repositories/ratelimitconfigurationrepo"
)

type App struct {
	notificationSrv ports.NotificationService
}

func NewApp() App {
	ns := notificationsrv.New(
		gateway.NewMailSender(),
		notificationrepo.NewInMemory(),
		ratelimitconfigurationrepo.NewInMemory(),
		utils.NewTimeProvider(),
	)

	return App{ns}
}

func (app App) SendNotification(userID string, notificationType string, message string) {
	err := app.notificationSrv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          message,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
