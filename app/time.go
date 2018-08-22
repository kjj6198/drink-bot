package app

import (
	"errors"
	"time"
)

// Time is a simple wrapper time.Time
type Time time.Time

const (
	// DefaultTimeZone is app timezone, defaults to Taipei
	DefaultTimeZone = "Asia/Taipei"
)

var (
	// ErrInvalidFormat emit if parse invalid time forma
	ErrInvalidFormat = errors.New("Invalid time format")
)

// Time returns Time in location
func (t Time) Time() time.Time {
	location, _ := time.LoadLocation(DefaultTimeZone)
	return time.Time(t).In(location)
}

// Parse returns time.Time in location
func (t Time) Parse(layout string, value string) (time.Time, error) {
	location, _ := time.LoadLocation(DefaultTimeZone)
	timeInLocation, err := time.ParseInLocation(layout, value, location)

	if err != nil {
		return time.Time{}, ErrInvalidFormat
	}

	return timeInLocation, nil
}
