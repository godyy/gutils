package worker

import (
	"log"
	"sync"
	"sync/atomic"
)

var debug = false

const (
	stateRunning  = 0
	stateStopping = 1
	stateStopped  = 2
)

// Workers routine池封装
type Workers struct {
	mtx         sync.RWMutex   // 锁
	state       atomic.Int32   // 状态
	idleWorkers chan *worker   // 空闲的routine
	workingWg   sync.WaitGroup // 工作中的routine等待
}

func NewWorkers(works int) *Workers {
	if works <= 0 {
		panic("works <= 0")
	}

	wp := &Workers{
		idleWorkers: make(chan *worker, works),
	}
	wp.state.Store(stateRunning)

	for i := 0; i < works; i++ {
		worker := newWorker(i)
		wp.idleWorkers <- worker
	}

	return wp
}

// Do 将任务交予空闲的routine执行
func (ws *Workers) Do(job func()) {
	if job == nil {
		return
	}

	ws.mtx.RLock()
	defer ws.mtx.RUnlock()

	if ws.state.Load() != stateRunning {
		return
	}

	worker := <-ws.idleWorkers
	ws.workingWg.Add(1)
	go worker.do(ws, job)
}

// Stop 停止所有routine的运行
func (ws *Workers) Stop() {
	ws.mtx.Lock()
	if ws.state.Load() != stateRunning {
		ws.mtx.Unlock()
		return
	}

	ws.state.Store(stateStopping)
	ws.mtx.Unlock()

	ws.workingWg.Wait()

	close(ws.idleWorkers)
	ws.state.Store(stateStopped)
}

func (ws *Workers) onIdleWorker(w *worker) {
	ws.idleWorkers <- w
	ws.workingWg.Done()
}

// worker 单routine封装
type worker struct {
	index int // 序号
}

func newWorker(index int) *worker {
	w := &worker{
		index: index,
	}
	return w
}

// do 将任务通过chan传递给routine执行
func (w *worker) do(ws *Workers, job func()) {
	if debug {
		log.Println("worker", w.index, "doing")
	}
	job()

	if debug {
		log.Println("worker", w.index, "idle")
	}
	ws.onIdleWorker(w)
}
