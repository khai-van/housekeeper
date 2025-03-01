package utils

import "time"

// parse string in format yyyy-mm-dd to time.Time
func ParseStandardDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
