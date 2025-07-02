package redis

import (
	"time"
	"unsafe"
)

const keyPrefix = "dankmuzikk:"

// duration returns duration based on the user's account id, where for regular user it's the set duration, and for guests it's 30 mins.
func duration(accountId uint64, duration time.Duration) time.Duration {
	if accountId > 1<<(8*unsafe.Sizeof(uint32(0))) {
		return time.Hour / 2
	}

	return duration
}
