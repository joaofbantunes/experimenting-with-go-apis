package shared

import "time"

type TimeProvider interface {
	Now() time.Time
}

type SystemTimeProvider struct{}

func NewSystemTimeProvider() SystemTimeProvider {
	return SystemTimeProvider{}
}

func (tp SystemTimeProvider) Now() time.Time {
	return time.Now()
}
