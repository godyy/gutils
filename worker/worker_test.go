package worker

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	debug = true
}

func TestWorkers(t *testing.T) {
	works := 1000
	ws := NewWorkers(works)

	var v atomic.Int32
	for i := 0; i < works; i++ {
		ws.Do(func() {
			time.Sleep(1 * time.Second)
			v.Add(1)
		})
	}

	ws.Stop()
	if v.Load() != int32(works) {
		t.Fatal()
	}

	ws = NewWorkers(works)
	wg := sync.WaitGroup{}
	n := 3
	v.Store(0)
	for i := 0; i < n; i++ {
		for j := 0; j < works; j++ {
			wg.Add(1)
			ws.Do(func() {
				time.Sleep(1 * time.Second)
				v.Add(1)
				wg.Done()
			})
		}
		wg.Wait()
		if v.Load() != int32((i+1)*works) {
			t.Fatal()
		}
	}
	ws.Stop()
}
