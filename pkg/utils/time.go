package utils

import "time"

func LocalTime() time.Time {
	return time.Now().UTC().Local()
}
