// Code generated by MockGen. DO NOT EDIT.
// Source: notification_service.go

// Package ports is a generated GoMock package.
package ports

import (
	reflect "reflect"

	domain "github.com/mendezdev/rate-limit-example/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockNotificationService is a mock of NotificationService interface.
type MockNotificationService struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceMockRecorder
}

// MockNotificationServiceMockRecorder is the mock recorder for MockNotificationService.
type MockNotificationServiceMockRecorder struct {
	mock *MockNotificationService
}

// NewMockNotificationService creates a new mock instance.
func NewMockNotificationService(ctrl *gomock.Controller) *MockNotificationService {
	mock := &MockNotificationService{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationService) EXPECT() *MockNotificationServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockNotificationService) Send(arg0 domain.NotificationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockNotificationServiceMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockNotificationService)(nil).Send), arg0)
}
