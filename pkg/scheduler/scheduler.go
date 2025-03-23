package scheduler

import (
	"container/heap"
	"sync"
	"sync/atomic"

	"github.com/gemirson/workpool/pkg/notification"
	"github.com/gemirson/workpool/pkg/task"
)

// noCopy is a placeholder type to prevent copying instances of sync.Mutex.
//
// The noCopy type is used to prevent copying instances of sync.Mutex.
// By embedding noCopy in a struct, the struct becomes uncopyable, which is useful
// when the struct contains a sync.Mutex field.
//
// This type is intended to be embedded in structs that require a mutex but do not
// need to be copied. By embedding noCopy, the struct becomes uncopyable, preventing
// accidental copying of the struct and its associated mutex.
type noCopy struct{}

// Lock is a placeholder method that does nothing.
//
// The Lock method is a placeholder method that does nothing.
// It is included to satisfy the sync.Locker interface, which requires a Lock method.
//
// This method is intended to be used in conjunction with the noCopy type to prevent
// copying instances of sync.Mutex. By embedding noCopy in a struct and implementing
// the Lock method as a no-op, the struct becomes uncopyable.
func (*noCopy) Lock() {}

// Unlock is a placeholder method that does nothing.
//
// The Unlock method is a placeholder method that does nothing.
// It is included to satisfy the sync.Locker interface, which requires an Unlock method.
//
// This method is intended to be used in conjunction with the noCopy type to prevent
// copying instances of sync.Mutex. By embedding noCopy in a struct and implementing
// the Unlock method as a no-op, the struct becomes uncopyable.
func (*noCopy) Unlock() {}

// priorityTask represents a task with a priority.
type priorityTask struct {
	task     *task.Task
	priority int
	index    int
}

// priorityQueue implements heap.Interface and holds priorityTasks.
type priorityQueue []*priorityTask

// Len returns the number of tasks in the priority queue.
//
// The Len method returns the number of tasks currently in the priority queue.
// It is used by the heap.Interface implementation to maintain the heap property.
//
// Return value:
// - int: The number of tasks in the priority queue.
func (pq priorityQueue) Len() int { return len(pq) }

// Less compares the priorities of two tasks in the priority queue.
// It returns true if the priority of the task at index i is greater than the priority of the task at index j.
// This function is used by the heap.Interface implementation to maintain the heap property in the priority queue.
//
// Parameters:
// - pq: A priorityQueue representing the priority queue.
// - i: An integer representing the index of the first task to compare.
// - j: An integer representing the index of the second task to compare.
//
// Return value:
// - bool: A boolean value indicating whether the priority of the task at index i is greater than the priority of the task at index j.
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

// Swap swaps the positions of two tasks in the priority queue based on their indices.
// It updates the indices of the swapped tasks to reflect their new positions in the queue.
//
// Parameters:
// - pq: A priorityQueue representing the priority queue.
// - i: An integer representing the index of the first task to swap.
// - j: An integer representing the index of the second task to swap.
//
// Return value:
// This function does not return any value.
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds a new priorityTask to the priority queue.
//
// The Push method adds a new priorityTask to the priority queue.
// It takes an interface x as a parameter, which is expected to be a pointer to a priorityTask.
// The method first calculates the current length of the priority queue (n) and then type-casts
// the interface x to a pointer to a priorityTask.
// It sets the index of the new priorityTask to the current length of the priority queue (n) and
// appends the new priorityTask to the end of the priority queue.
//
// Parameters:
//   - x: An interface representing the new priorityTask to be added to the priority queue.
//     It is expected to be a pointer to a priorityTask.
//
// Return value:
// This function does not return any value.
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*priorityTask)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the highest priority task from the priority queue.
//
// The Pop method removes and returns the highest priority task from the priority queue.
// It is used by the heap.Interface implementation to maintain the heap property.
//
// The method first copies the current priority queue into a slice called old.
// It then calculates the length of old (n) and retrieves the last element in old,
// which represents the highest priority task.
// The last element in old is set to nil and its index is updated to -1 to indicate
// that it is no longer in the priority queue.
// Finally, the priority queue is updated by slicing old from index 0 to n-1,
// effectively removing the highest priority task from the queue.
//
// Parameters:
// This method does not take any parameters.
//
// Return value:
//   - interface{}: The interface{} type represents any Go value.
//     In this case, the method returns the highest priority task as an interface{}.
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Workpool represents a pool of workers that execute tasks concurrently.
// It manages a queue of tasks with priorities and assigns tasks to available workers.
// The Workpool ensures that tasks are executed in the correct order based on their priorities.
type Workpool struct {
	noCopy      noCopy
	workerCount int
	taskQueue   priorityQueue
	mu          sync.Mutex
	cond        *sync.Cond
	taskMap     map[string]*priorityTask
	notifier    *notification.Notifier
	running     int32
	// nextWorker is an integer that keeps track of the index of the next worker
	// that should be assigned a task in the work pool. It is used to distribute
	// tasks in a round-robin fashion among available workers.
	nextWorker int
	workers    []chan *task.Task
}

// NewWorkpool creates a new instance of Workpool with the specified number of workers.
// The Workpool manages a queue of tasks with priorities and assigns tasks to available workers.
// It ensures that tasks are executed in the correct order based on their priorities.
//
// Parameters:
// - workerCount: An integer representing the number of workers to be created in the workpool.
//
// Return value:
// - *Workpool: A pointer to the newly created Workpool instance.
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
// Submit adds a task to the workpool with the given priority.
//
// The function Submit adds a new task to the workpool with the specified priority.
// It locks the workpool's mutex to ensure thread safety while modifying the task queue and task map.
// The new task is wrapped in a priorityTask struct and added to the task queue using heap.Push.
// The task map is updated with the task's ID and priorityTask pointer.
// Finally, the condition variable is signaled to wake up a waiting worker to process the new task.
//
// Parameters:
// - t: A pointer to the task.Task struct representing the task to be added to the workpool.
// - priority: An integer representing the priority of the task. Higher values indicate higher priority.
//
// Return value:
// This function does not return any value.
func (wp *Workpool) Submit(t *task.Task, priority int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	pt := &priorityTask{task: t, priority: priority}
	heap.Push(&wp.taskQueue, pt)
	wp.taskMap[t.ID()] = pt
	wp.cond.Signal()
}

// worker processes tasks from the task channel.
// worker processes tasks from the task channel.
// It continuously monitors the task queue for new tasks and assigns them to available workers.
// Once a task is assigned, the worker executes the task, notifies the notifier about the task's completion,
// and sends the task back to the task channel.
// The worker continues to process tasks until the workpool is shut down.
func (wp *Workpool) worker(taskChan chan *task.Task) {
	for {
		wp.mu.Lock()
		// Wait for a task to be available in the task queue.
		for wp.taskQueue.Len() == 0 {
			wp.cond.Wait()
		}

		// Pop the highest priority task from the task queue.
		pt := heap.Pop(&wp.taskQueue).(*priorityTask)

		// Remove the task from the task map.
		delete(wp.taskMap, pt.task.ID())
		wp.mu.Unlock()

		// Increment the running counter to indicate that a task is being processed.
		atomic.AddInt32(&wp.running, 1)

		// Execute the task.
		pt.task.Execute()

		// Notify the notifier about the task's completion.
		wp.notifier.Notify(pt.task.ID())

		// Send the task back to the task channel.
		taskChan <- pt.task

		// Decrement the running counter to indicate that a task has been completed.
		atomic.AddInt32(&wp.running, -1)
	}
}

// removeCompletedTasks removes tasks with state Completed from the task queue.
// removeCompletedTasks removes tasks with state Completed from the task queue.
// It iterates through the task map, checks the state of each task, and removes
// the task from the task queue and the task map if its state is Completed.
// This function is called when shutting down the workpool to ensure that all
// completed tasks are removed from the task queue and task map.
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

// Cancel cancels a task with the given taskID in the workpool.
// It locks the workpool's mutex to ensure thread safety while modifying the task queue and task map.
// If the task with the given taskID exists in the task map, it removes the task from the task queue,
// deletes the task from the task map, and cancels the task.
//
// Parameters:
// - taskID: A string representing the ID of the task to be canceled.
//
// Return value:
// This function does not return any value.
func (wp *Workpool) Cancel(taskID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if pt, ok := wp.taskMap[taskID]; ok {
		heap.Remove(&wp.taskQueue, pt.index)
		delete(wp.taskMap, taskID)
		pt.task.Cancel()
	}
}

// Shutdown gracefully shuts down the workpool.
// Shutdown gracefully shuts down the workpool.
// It ensures that all tasks are canceled, closes the worker channels, waits for all tasks to complete,
// and removes completed tasks from the task queue and task map.
func (wp *Workpool) Shutdown() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	// Cancel all tasks in the task map.
	for _, pt := range wp.taskMap {
		wp.Cancel(pt.task.ID())
	}

	// Close all worker channels.
	for _, worker := range wp.workers {
		close(worker)
	}

	// Wait for all tasks to complete.
	for atomic.LoadInt32(&wp.running) > 0 {
		wp.cond.Wait()
	}

	// Remove completed tasks from the task queue and task map.
	wp.removeCompletedTasks()
}
