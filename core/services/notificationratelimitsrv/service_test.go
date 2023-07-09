package notificationratelimitsrv

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

func TestShouldBeExceededWhenRateLimitIsExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockRateLimitConfigRepo, mockNotificationRepo, mockTimeProvider)

	userID := "userID"
	notificationType := "notification_type"

	oneSecondInThePast := timeNow.Add(-time.Duration(1) * time.Second)
	oneMinuteInThePast := timeNow.Add(-time.Duration(1) * time.Minute)
	oneHourInThePast := timeNow.Add(-time.Duration(1) * time.Hour)
	fiveHourInThePast := timeNow.Add(-time.Duration(5) * time.Hour)

	testCases := []struct {
		Name            string
		Notifications   []domain.Notification
		RateLimitConfig domain.RateLimitConfiguration
		DateFrom        time.Time
	}{
		{
			Name:          "1 notification per second case",
			Notifications: createNotificationsMock(userID, notificationType, oneSecondInThePast, 1),
			RateLimitConfig: domain.RateLimitConfiguration{
				Name:        notificationType,
				Limit:       1,
				TimeUnit:    1,
				TimeMeasure: "SECONDS",
			},
			DateFrom: oneSecondInThePast,
		},
		{
			Name:          "1 notifications per minute case",
			Notifications: createNotificationsMock(userID, notificationType, oneMinuteInThePast, 3),
			RateLimitConfig: domain.RateLimitConfiguration{
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
			RateLimitConfig: domain.RateLimitConfiguration{
				Name:        notificationType,
				Limit:       10,
				TimeUnit:    1,
				TimeMeasure: "HOURS",
			},
			DateFrom: oneHourInThePast,
		},
		{
			Name:          "20 notifications per 5 hour case",
			Notifications: createNotificationsMock(userID, notificationType, fiveHourInThePast, 20),
			RateLimitConfig: domain.RateLimitConfiguration{
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
				Return(&tc.RateLimitConfig, nil).
				Times(1)

			mockNotificationRepo.
				EXPECT().
				GetByTypeAndUserAndFromDate(userID, notificationType, tc.DateFrom).
				Return(tc.Notifications, nil).
				Times(1)

			result := srv.IsExceeded(userID, notificationType)
			assert.True(t, result)
		})
	}
}

func TestShouldNotBeExceededWhenRateLimitIsNotExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockRateLimitConfigRepo, mockNotificationRepo, mockTimeProvider)

	userID := "userID"
	notificationType := "notification_type"

	oneSecondInThePast := timeNow.Add(-time.Duration(1) * time.Second)
	notifications := createNotificationsMock(userID, notificationType, oneSecondInThePast, 1)
	rlc := domain.RateLimitConfiguration{
		Name:        notificationType,
		Limit:       2,
		TimeUnit:    1,
		TimeMeasure: "SECONDS",
	}

	mockTimeProvider.
		EXPECT().
		Now().
		Return(timeNow).
		AnyTimes()

	mockRateLimitConfigRepo.
		EXPECT().
		GetRateLimitConfiguration(notificationType).
		Return(&rlc, nil).
		Times(1)

	mockNotificationRepo.
		EXPECT().
		GetByTypeAndUserAndFromDate(userID, notificationType, oneSecondInThePast).
		Return(notifications, nil).
		Times(1)

	result := srv.IsExceeded(userID, notificationType)
	assert.False(t, result)
}

func TestShouldNotBeExceededWhenRateLimitConfigRepoReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockRateLimitConfigRepo, mockNotificationRepo, mockTimeProvider)

	userID := "userID"
	notificationType := "notification_type"

	mockRateLimitConfigRepo.
		EXPECT().
		GetRateLimitConfiguration(notificationType).
		Return(nil, errors.New("db_error")).
		Times(1)

	result := srv.IsExceeded(userID, notificationType)
	assert.False(t, result)
}

func TestShouldNotBeExceededWhenRateLimitConfigIsNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRateLimitConfigRepo := mocks.NewMockRateLimitConfigurationRepository(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockRateLimitConfigRepo, mockNotificationRepo, mockTimeProvider)

	userID := "userID"
	notificationType := "notification_type"

	mockRateLimitConfigRepo.
		EXPECT().
		GetRateLimitConfiguration(notificationType).
		Return(nil, nil).
		Times(1)

	result := srv.IsExceeded(userID, notificationType)
	assert.False(t, result)
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
