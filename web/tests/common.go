package tests

import (
	"math/rand"
	"time"
)

var (
	random *rand.Rand
)

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixMilli()))
	initAccounts()
	initProfiles()
	initPlaylists()
}
