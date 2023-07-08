package notificationsrv

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/mendezdev/rate-limit-example/core/domain"
	mocks "github.com/mendezdev/rate-limit-example/core/ports/mocks"
	"github.com/mendezdev/rate-limit-example/core/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	timeNow = time.Date(2018, time.January, 1, 10, 0, 0, 0, time.UTC)
)

func TestShouldSendMailNotificationWhenRateLimitIsNotExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRepo, mockRateLimitConfigRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"

	mockTimeProvider.
		EXPECT().
		Now().
		Return(timeNow).
		AnyTimes()

	rlc := domain.RateLimitConfiguration{
		Name:        notificationType,
		Limit:       3,
		TimeUnit:    1,
		TimeMeasure: "MINUTES",
	}
	mockRateLimitConfigRepo.
		EXPECT().
		GetRateLimitConfiguration(notificationType).
		Return(&rlc, nil).
		Times(1)

	notifications := make([]domain.Notification, 0)

	mockNotificationRepo.
		EXPECT().
		GetByTypeAndUserAndFromDate(userID, notificationType, timeNow.Add(-time.Duration(1)*time.Minute)).
		Return(notifications, nil).
		Times(1)

	newNotificationMsg := "some message 3"
	mockNotificationRepo.
		EXPECT().
		Save(domain.Notification{
			UserID:           userID,
			NotificationType: notificationType,
			DateCreated:      timeNow,
			Message:          newNotificationMsg,
		})

	mockMailSender.
		EXPECT().
		Send(userID, newNotificationMsg).
		Return(nil).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.Nil(t, err)
}

func createNotificationsMock(
	userID string,
	notificationType string,
	dateCreated time.Time,
	notificationQty int,
) []domain.Notification {
	notifications := make([]domain.Notification, 0)
	for i := 1; i <= notificationQty; i++ {
		notifications = append(notifications, domain.Notification{
			UserID:           userID,
			NotificationType: notificationType,
			DateCreated:      dateCreated,
			Message:          fmt.Sprintf("message %d", i),
		})
	}

	return notifications
}

func TestShouldNotSendMailNotificationWhenRateLimitIsExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRepo, mockRateLimitConfigRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"
	notificationRequest := domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          "new message",
	}
	oneSecondInThePast := timeNow.Add(-time.Duration(1) * time.Second)
	oneMinuteInThePast := timeNow.Add(-time.Duration(1) * time.Minute)
	oneHourInThePast := timeNow.Add(-time.Duration(1) * time.Hour)
	fiveHourInThePast := timeNow.Add(-time.Duration(5) * time.Hour)

	testCases := []struct {
		Name          string
		Notifications []domain.Notification
		RateLimitCfg  domain.RateLimitConfiguration
		DateFrom      time.Time
	}{
		{
			Name:          "1 notification per second case",
			Notifications: createNotificationsMock(userID, notificationType, oneSecondInThePast, 1),
			RateLimitCfg: domain.RateLimitConfiguration{
				Name:        notificationType,
				Limit:       1,
				TimeUnit:    1,
				TimeMeasure: "SECONDS",
			},
			DateFrom: oneSecondInThePast,
		},
		{
			Name:          "3 notifications per minute case",
			Notifications: createNotificationsMock(userID, notificationType, oneMinuteInThePast, 3),
			RateLimitCfg: domain.RateLimitConfiguration{
				Name:        notificationType,
				Limit:       3,
				TimeUnit:    1,
				TimeMeasure: "MINUTES",
			},
			DateFrom: oneMinuteInThePast,
		},
		{
			Name:          "10 notifications per hour case",
			Notifications: createNotificationsMock(userID, notificationType, oneHourInThePast, 10),
			RateLimitCfg: domain.RateLimitConfiguration{
				Name:        notificationType,
				Limit:       10,
				TimeUnit:    1,
				TimeMeasure: "HOURS",
			},
			DateFrom: oneHourInThePast,
		},
		{
			Name:          "20 notifications per 5 hours case",
			Notifications: createNotificationsMock(userID, notificationType, fiveHourInThePast, 20),
			RateLimitCfg: domain.RateLimitConfiguration{
				Name:        notificationType,
				Limit:       20,
				TimeUnit:    5,
				TimeMeasure: "HOURS",
			},
			DateFrom: fiveHourInThePast,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockTimeProvider.
				EXPECT().
				Now().
				Return(timeNow).
				AnyTimes()

			mockRateLimitConfigRepo.
				EXPECT().
				GetRateLimitConfiguration(notificationType).
				Return(&tc.RateLimitCfg, nil).
				Times(1)

			mockNotificationRepo.
				EXPECT().
				GetByTypeAndUserAndFromDate(userID, notificationType, tc.DateFrom).
				Return(tc.Notifications, nil).
				Times(1)

			err := srv.Send(notificationRequest)

			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), "the rate limit has been exceeded for user_id "+userID+" and notification type "+notificationType)
		})
	}
}

func TestShouldSendMailNotificationWhenNoRateLimitConfiguration(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRepo, mockRateLimitConfigRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"

	mockTimeProvider.
		EXPECT().
		Now().
		Return(timeNow).
		AnyTimes()

	mockRateLimitConfigRepo.
		EXPECT().
		GetRateLimitConfiguration(notificationType).
		Return(nil, nil).
		Times(1)

	newNotificationMsg := "new message"
	mockNotificationRepo.
		EXPECT().
		Save(domain.Notification{
			UserID:           userID,
			NotificationType: notificationType,
			DateCreated:      timeNow,
			Message:          newNotificationMsg,
		})

	mockMailSender.
		EXPECT().
		Send(userID, newNotificationMsg).
		Return(nil).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.Nil(t, err)
}

func TestShouldSendMailNotificationWhenGetRateLimitConfigurationReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRepo, mockRateLimitConfigRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"

	mockTimeProvider.
		EXPECT().
		Now().
		Return(timeNow).
		AnyTimes()

	mockRateLimitConfigRepo.
		EXPECT().
		GetRateLimitConfiguration(notificationType).
		Return(nil, errors.New("rate limit configuration repo error")).
		Times(1)

	newNotificationMsg := "new message"
	mockNotificationRepo.
		EXPECT().
		Save(domain.Notification{
			UserID:           userID,
			NotificationType: notificationType,
			DateCreated:      timeNow,
			Message:          newNotificationMsg,
		})

	mockMailSender.
		EXPECT().
		Send(userID, newNotificationMsg).
		Return(nil).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.Nil(t, err)
}
