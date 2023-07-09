package notificationsrv

import (
	"errors"
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

func TestShouldSendMailNotification(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRateLimitSrv := mocks.NewMockNotificationRateLimitService(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRateLimitSrv, mockNotificationRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"
	newNotificationMsg := "new message"

	mockNotificationRateLimitSrv.
		EXPECT().
		IsExceeded(userID, notificationType).
		Return(false).
		Times(1)

	mockMailSender.
		EXPECT().
		Send(userID, newNotificationMsg).
		Return(nil).
		Times(1)

	mockTimeProvider.
		EXPECT().
		Now().
		Return(timeNow).
		AnyTimes()

	mockNotificationRepo.
		EXPECT().
		Save(domain.Notification{
			UserID:           userID,
			NotificationType: notificationType,
			DateCreated:      timeNow,
			Message:          newNotificationMsg,
		}).
		Return(nil).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.Nil(t, err)
}

func TestShouldNotSendMailNotificationWhenRateLimitIsExceeded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRateLimitSrv := mocks.NewMockNotificationRateLimitService(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRateLimitSrv, mockNotificationRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"
	newNotificationMsg := "new message"

	mockNotificationRateLimitSrv.
		EXPECT().
		IsExceeded(userID, notificationType).
		Return(true).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "the rate limit has been exceeded for user_id "+userID+" and notification type "+notificationType)
}

func TestShouldNotSendMailNotificationWhenMailSenderReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRateLimitSrv := mocks.NewMockNotificationRateLimitService(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRateLimitSrv, mockNotificationRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"
	newNotificationMsg := "new message"

	mockNotificationRateLimitSrv.
		EXPECT().
		IsExceeded(userID, notificationType).
		Return(false).
		Times(1)

	mockMailSender.
		EXPECT().
		Send(userID, newNotificationMsg).
		Return(errors.New("mailsender_error")).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error trying to send email: mailsender_error")
}

func TestShouldSendMailNotificationWhenSaveNotificationReturnsError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailSender := mocks.NewMockMailSenderGateway(mockCtrl)
	mockNotificationRateLimitSrv := mocks.NewMockNotificationRateLimitService(mockCtrl)
	mockNotificationRepo := mocks.NewMockNotificationRepository(mockCtrl)
	mockTimeProvider := utils.NewMockTimeProvider(mockCtrl)

	srv := New(mockMailSender, mockNotificationRateLimitSrv, mockNotificationRepo, mockTimeProvider)

	notificationType := "notification_type"
	userID := "userID"
	newNotificationMsg := "new message"

	mockNotificationRateLimitSrv.
		EXPECT().
		IsExceeded(userID, notificationType).
		Return(false).
		Times(1)

	mockMailSender.
		EXPECT().
		Send(userID, newNotificationMsg).
		Return(nil).
		Times(1)

	mockTimeProvider.
		EXPECT().
		Now().
		Return(timeNow).
		AnyTimes()

	mockNotificationRepo.
		EXPECT().
		Save(domain.Notification{
			UserID:           userID,
			NotificationType: notificationType,
			DateCreated:      timeNow,
			Message:          newNotificationMsg,
		}).
		Return(errors.New("save_notification_error")).
		Times(1)

	err := srv.Send(domain.NotificationRequest{
		UserID:           userID,
		NotificationType: notificationType,
		Message:          newNotificationMsg,
	})

	assert.Nil(t, err)
}
