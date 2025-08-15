package util

import "time"

func ClearYear(t time.Time) time.Time {
	return time.Date(
		0,
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		t.Location(),
	)
}
