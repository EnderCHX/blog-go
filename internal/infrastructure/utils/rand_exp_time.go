package utils

import (
	"math/rand"
	"time"
)

// from start to end minutes
func RandExpTime(start, end int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return time.Minute * time.Duration(start+r.Intn(end-start)+1)
}
