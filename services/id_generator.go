package services

import "sync/atomic"

func GenerateID(counter *int64) int64 {
	return atomic.AddInt64(counter, 1)
}
