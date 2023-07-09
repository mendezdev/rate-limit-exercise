package notificationrepo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mendezdev/rate-limit-example/core/domain"
	"github.com/mendezdev/rate-limit-example/core/ports"
)

const (
	path = "./mock_store/notifications.json"
)

var (
	byNotificationTypeInitiazer map[string][]domain.Notification
)

type inmemory struct {
	// we could have another nested map by userID
	byNotificationType map[string][]domain.Notification
}

// GetByTypeAndUserAndFromDate implements ports.NotificationRepository.
func (mem *inmemory) GetByTypeAndUserAndFromDate(userID string, notificationType string, fromDate time.Time) ([]domain.Notification, error) {
	notifications, exist := mem.byNotificationType[notificationType]
	if !exist {
		return nil, nil
	}

	notificationsFiltered := make([]domain.Notification, 0)
	for _, n := range notifications {
		if n.UserID == userID && n.DateCreated.Compare(fromDate) >= 0 {
			notificationsFiltered = append(notificationsFiltered, n)
		}
	}

	return notificationsFiltered, nil
}

func NewInMemory() ports.NotificationRepository {
	return &inmemory{
		byNotificationType: byNotificationTypeInitiazer,
	}
}

// Save implements ports.NotificationRepository.
func (mem *inmemory) Save(n domain.Notification) error {
	_, ok := mem.byNotificationType[n.NotificationType]
	if !ok {
		mem.byNotificationType[n.NotificationType] = make([]domain.Notification, 0)
	}
	mem.byNotificationType[n.NotificationType] = append(mem.byNotificationType[n.NotificationType], n)
	return nil
}

func init() {
	fmt.Println("initializing inmemory notification db...")

	byNotificationTypeInitiazer = make(map[string][]domain.Notification)

	var notifications []domain.Notification
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("error reading json file for initialize notification db: %s", err.Error()))
	}

	jsonErr := json.Unmarshal(data, &notifications)
	if jsonErr != nil {
		panic(fmt.Errorf("error unmarshaling json data: %s", jsonErr.Error()))
	}

	for _, n := range notifications {
		_, ok := byNotificationTypeInitiazer[n.NotificationType]
		if !ok {
			byNotificationTypeInitiazer[n.NotificationType] = make([]domain.Notification, 0)
		}
		byNotificationTypeInitiazer[n.NotificationType] = append(byNotificationTypeInitiazer[n.NotificationType], n)
	}
	fmt.Println("inmemory notification db initilized.")
}
