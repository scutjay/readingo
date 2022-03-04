package util

import "time"

// ParseStringToDuration return 1 min as default
func ParseStringToDuration(str string) time.Duration {
	duration, err := time.ParseDuration(str)
	if err != nil {
		return time.Minute
	}
	return duration
}
