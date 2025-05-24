package utils

import "time"

func FormatedTime(t time.Time) string {
	// UTC+0
	return t.UTC().Format("2006-01-02 15:04:05")
}
