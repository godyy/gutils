package leaky_bucket

import (
	"log"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	limiter := NewMutexBased(10, WithPer(time.Second))
	limiter.Take()
	time.Sleep(2 * time.Second)
	for i := 0; i < 100; i++ {
		limiter.Take()
		log.Printf("%d pass", i)
	}
}

func TestAtomicInt64(t *testing.T) {
	limiter := NewAtomicInt64Based(10, WithPer(time.Second))
	limiter.Take()
	time.Sleep(2 * time.Second)
	for i := 0; i < 100; i++ {
		limiter.Take()
		log.Printf("%d pass", i)
	}
}
