package utils

import "time"

//go:generate mockgen -destination=mock_time_provider.go -package=utils -source=time_provider.go TimeProvider

type TimeProvider interface {
	Now() time.Time
}

type realTimeProvider struct{}

// Now implements TimeProvider.
func (rtp realTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

func NewTimeProvider() TimeProvider {
	return realTimeProvider{}
}
