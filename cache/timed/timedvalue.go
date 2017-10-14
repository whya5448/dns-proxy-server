package timed

import (
	"time"
)

type TimedValue interface {
	Creation() time.Time
	Timeout() time.Duration
	Value() interface{}

	//
	// Check if the value has expired
	// now current time to be compared with the Creation()
	//
	IsValid(now time.Time) bool

}

type timedValueImpl struct {
	creation time.Time
	timeout time.Duration
	value interface{}
}

func(t *timedValueImpl) Creation() time.Time {
	return t.creation
}

func(t *timedValueImpl) Timeout() time.Duration {
	return t.timeout
}

func(t *timedValueImpl) Value() interface{} {
	return t.value
}

func(t *timedValueImpl) IsValid(now time.Time) bool {
	return t.Timeout() > now.Sub(t.Creation())
}

func NewTimedValue(value interface{}, creation time.Time, timeout time.Duration) TimedValue {
	return &timedValueImpl{creation:creation, value:value, timeout:timeout}
}
