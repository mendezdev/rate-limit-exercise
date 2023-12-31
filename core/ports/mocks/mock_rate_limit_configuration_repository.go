// Code generated by MockGen. DO NOT EDIT.
// Source: rate_limit_configuration_repository.go

// Package ports is a generated GoMock package.
package ports

import (
	reflect "reflect"

	domain "github.com/mendezdev/rate-limit-example/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRateLimitConfigurationRepository is a mock of RateLimitConfigurationRepository interface.
type MockRateLimitConfigurationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRateLimitConfigurationRepositoryMockRecorder
}

// MockRateLimitConfigurationRepositoryMockRecorder is the mock recorder for MockRateLimitConfigurationRepository.
type MockRateLimitConfigurationRepositoryMockRecorder struct {
	mock *MockRateLimitConfigurationRepository
}

// NewMockRateLimitConfigurationRepository creates a new mock instance.
func NewMockRateLimitConfigurationRepository(ctrl *gomock.Controller) *MockRateLimitConfigurationRepository {
	mock := &MockRateLimitConfigurationRepository{ctrl: ctrl}
	mock.recorder = &MockRateLimitConfigurationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateLimitConfigurationRepository) EXPECT() *MockRateLimitConfigurationRepositoryMockRecorder {
	return m.recorder
}

// GetRateLimitConfiguration mocks base method.
func (m *MockRateLimitConfigurationRepository) GetRateLimitConfiguration(configType string) (*domain.RateLimitConfiguration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRateLimitConfiguration", configType)
	ret0, _ := ret[0].(*domain.RateLimitConfiguration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRateLimitConfiguration indicates an expected call of GetRateLimitConfiguration.
func (mr *MockRateLimitConfigurationRepositoryMockRecorder) GetRateLimitConfiguration(configType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRateLimitConfiguration", reflect.TypeOf((*MockRateLimitConfigurationRepository)(nil).GetRateLimitConfiguration), configType)
}
