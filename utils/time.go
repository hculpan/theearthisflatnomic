package utils

import "time"

var location *time.Location

// GetLocation returns the locations for all
// time functions
func GetLocation() *time.Location {
	if location == nil {
		location, _ = time.LoadLocation("EST")
	}
	return location
}

// GetEndTime returns the time days from
// now at 11:59pm
func GetEndTime(days int) time.Time {
	var result time.Time

	loc := GetLocation()
	n := time.Now().In(loc)
	result = time.Date(n.Year(), n.Month(), n.Day(), 23, 59, 0, 0, loc)
	result = result.Add((4 * 24) * time.Hour)

	return result
}
