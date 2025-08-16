package util

import "time"

func ClearYear(t time.Time) time.Time {
	return time.Date(
		0,
		t.Month(),
		t.Day(),
		0,
		0,
		0,
		0,
		t.Location(),
	)
}
