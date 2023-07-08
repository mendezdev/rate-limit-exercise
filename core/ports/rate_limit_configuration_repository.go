package ports

import "github.com/mendezdev/rate-limit-example/core/domain"

//go:generate mockgen -destination=mocks/mock_rate_limit_configuration_repository.go -package=ports -source=rate_limit_configuration_repository.go RateLimitConfigurationRepository

type RateLimitConfigurationRepository interface {
	GetRateLimitConfiguration(configType string) (*domain.RateLimitConfiguration, error)
}
