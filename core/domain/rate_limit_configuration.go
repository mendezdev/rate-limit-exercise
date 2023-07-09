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

func (rlc RateLimitConfiguration) GetTimeMeasureInDuration() time.Duration {
	switch rlc.TimeMeasure {
	case "SECONDS":
		return time.Second
	case "MINUTES":
		return time.Minute
	case "HOURS":
		return time.Hour
	default:
		// TODO: maybe we return an error instead of throw a panic and shut down the program
		panic(fmt.Sprintf("no time measure found for %s", rlc.TimeMeasure))
	}
}
