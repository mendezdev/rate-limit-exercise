// Code generated by MockGen. DO NOT EDIT.
// Source: notification_repository.go

// Package ports is a generated GoMock package.
package ports

import (
	reflect "reflect"
	time "time"

	domain "github.com/mendezdev/rate-limit-example/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockNotificationRepository is a mock of NotificationRepository interface.
type MockNotificationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationRepositoryMockRecorder
}

// MockNotificationRepositoryMockRecorder is the mock recorder for MockNotificationRepository.
type MockNotificationRepositoryMockRecorder struct {
	mock *MockNotificationRepository
}

// NewMockNotificationRepository creates a new mock instance.
func NewMockNotificationRepository(ctrl *gomock.Controller) *MockNotificationRepository {
	mock := &MockNotificationRepository{ctrl: ctrl}
	mock.recorder = &MockNotificationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationRepository) EXPECT() *MockNotificationRepositoryMockRecorder {
	return m.recorder
}

// GetByTypeAndUserAndFromDate mocks base method.
func (m *MockNotificationRepository) GetByTypeAndUserAndFromDate(userID, notificationType string, fromDate time.Time) ([]domain.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTypeAndUserAndFromDate", userID, notificationType, fromDate)
	ret0, _ := ret[0].([]domain.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTypeAndUserAndFromDate indicates an expected call of GetByTypeAndUserAndFromDate.
func (mr *MockNotificationRepositoryMockRecorder) GetByTypeAndUserAndFromDate(userID, notificationType, fromDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTypeAndUserAndFromDate", reflect.TypeOf((*MockNotificationRepository)(nil).GetByTypeAndUserAndFromDate), userID, notificationType, fromDate)
}

// Save mocks base method.
func (m *MockNotificationRepository) Save(arg0 domain.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockNotificationRepositoryMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockNotificationRepository)(nil).Save), arg0)
}