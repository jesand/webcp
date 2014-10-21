package main

import (
	"errors"
	"time"
)

func ParseDate(date string) (time.Time, error) {
	switch len(date) {
	case 4:
		return time.Parse("2006", date)
	case 6:
		return time.Parse("200601", date)
	case 8:
		return time.Parse("20060102", date)
	case 14:
		return time.Parse("20060102150405", date)
	default:
		return time.Time{}, errors.New("Invalid date format.")
	}
}
