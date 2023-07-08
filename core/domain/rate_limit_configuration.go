package domain

import (
	"fmt"
	"time"
)

type RateLimitConfiguration struct {
	Name        string `json:"name"`
	Limit       int    `json:"limit"`
	TimeUnit    int    `json:"time_unit"`
	TimeMeasure string `json:"time_measure"`
}

func (rlc RateLimitConfiguration) getTimeMeasureInDuration() time.Duration {
	switch rlc.TimeMeasure {
	case "SECONDS":
		return time.Second
	case "MINUTES":
		return time.Minute
	case "HOURS":
		return time.Hour
	default:
		// TODO: review this
		panic(fmt.Sprintf("no time measure found for %s", rlc.TimeMeasure))
	}
}

func (rlc RateLimitConfiguration) GetDateFromFor(date time.Time) time.Time {
	return date.Add(-time.Duration(rlc.TimeUnit) * rlc.getTimeMeasureInDuration())
}