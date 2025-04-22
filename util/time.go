package util

import "time"

func TimeExpired(t time.Time) bool {
	now := time.Now()
	expired := now.After(t)
	return expired
}
