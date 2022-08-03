package util

import "time"

const (
	layout = "2006-01-02 15:04:05"
)

func Now() string {
	return time.Now().Format(layout)
}

func String2Time(s string) time.Time {
	t, _ := time.ParseInLocation(layout, s, time.Local)
	return t
}
