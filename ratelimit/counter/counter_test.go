package counter

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	c := New(1000, time.Second)
	for i := 0; i < 1010; i++ {
		if i < 1000 {
			if !c.Allow() {
				t.Fatal("allow should be true")
			}
		} else {
			if c.Allow() {
				t.Fatal("allow should be false")
			}
		}
	}

	time.Sleep(time.Second)

	if !c.Allow() {
		t.Fatal("allow should be true last")
	}
}
