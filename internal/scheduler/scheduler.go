package scheduler

import (
	"container/heap"
	"sync"
	"sync/atomic"

	"github.com/gemirson/workpool/internal/notification"
	"github.com/gemirson/workpool/internal/task"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// priorityTask represents a task with a priority.
type priorityTask struct {
	task     *task.Task
	priority int
	index    int
}

// priorityQueue implements heap.Interface and holds priorityTasks.
type priorityQueue []*priorityTask

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*priorityTask)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Workpool manages a pool of workers to execute tasks.
type Workpool struct {
	noCopy      noCopy
	workerCount int
	taskQueue   priorityQueue
	mu          sync.Mutex
	cond        *sync.Cond
	taskMap     map[string]*priorityTask
	notifier    *notification.Notifier
	running     int32
	nextWorker  int
	workers     []chan *task.Task
}

// NewWorkpool creates a new Workpool with the given number of workers.
func NewWorkpool(workerCount int) *Workpool {
	wp := &Workpool{
		workerCount: workerCount,
		taskQueue:   make(priorityQueue, 0),
		taskMap:     make(map[string]*priorityTask),
		notifier:    notification.NewNotifier(),
		workers:     make([]chan *task.Task, workerCount),
	}
	wp.cond = sync.NewCond(&wp.mu)
	for i := 0; i < workerCount; i++ {
		wp.workers[i] = make(chan *task.Task)
		go wp.worker(wp.workers[i])
	}
	return wp
}

// Submit adds a task to the workpool with the given priority.
func (wp *Workpool) Submit(t *task.Task, priority int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	pt := &priorityTask{task: t, priority: priority}
	heap.Push(&wp.taskQueue, pt)
	wp.taskMap[t.ID()] = pt
	wp.cond.Signal()
}

// Cancel removes a task from the workpool by its ID.
func (wp *Workpool) Cancel(taskID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if pt, ok := wp.taskMap[taskID]; ok {
		heap.Remove(&wp.taskQueue, pt.index)
		delete(wp.taskMap, taskID)
		pt.task.Cancel()
	}
}

// worker processes tasks from the task channel.
func (wp *Workpool) worker(taskChan chan *task.Task) {
	for {
		wp.mu.Lock()
		for wp.taskQueue.Len() == 0 {
			wp.cond.Wait()
		}
		pt := heap.Pop(&wp.taskQueue).(*priorityTask)
		delete(wp.taskMap, pt.task.ID())
		wp.mu.Unlock()
		atomic.AddInt32(&wp.running, 1)
		taskChan <- pt.task
		pt.task.Execute()
		wp.notifier.Notify(pt.task.ID())
		atomic.AddInt32(&wp.running, -1)
	}
}

// removeCompletedTasks removes tasks with state Completed from the task queue.
func (wp *Workpool) removeCompletedTasks() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for id, pt := range wp.taskMap {
		if pt.task.State() == task.Completed {
			heap.Remove(&wp.taskQueue, pt.index)
			delete(wp.taskMap, id)
		}
	}
}

// Shutdown gracefully shuts down the workpool.
func (wp *Workpool) Shutdown() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for atomic.LoadInt32(&wp.running) > 0 {
		wp.cond.Wait()
	}
	for _, pt := range wp.taskMap {
		pt.task.Cancel()
	}
	for _, worker := range wp.workers {
		close(worker)
	}
	wp.removeCompletedTasks()
}
