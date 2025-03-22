package task

import (
	"sync/atomic"

	"github.com/google/uuid"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// TaskState represents the state of a task.
type TaskState int32

const (
	Pending TaskState = iota
	Running
	Completed
	Cancelled
)

// Task represents a unit of work to be executed.
type Task struct {
	noCopy     noCopy
	id         string
	action     func()
	cancel     int32
	cancelChan chan struct{}
	state      TaskState
}

// NewTask creates a new Task with the given action.
func NewTask(action func()) *Task {
	return &Task{
		id:         uuid.New().String(),
		action:     action,
		cancel:     0,
		cancelChan: make(chan struct{}),
		state:      Pending,
	}
}

// Execute runs the task's action.
func (t *Task) Execute() {
	if atomic.LoadInt32(&t.cancel) == 1 {
		atomic.StoreInt32((*int32)(&t.state), int32(Cancelled))
		return
	}
	atomic.StoreInt32((*int32)(&t.state), int32(Running))
	t.action()
	atomic.StoreInt32((*int32)(&t.state), int32(Completed))
}

// Cancel cancels the task.
func (t *Task) Cancel() {
	if atomic.CompareAndSwapInt32(&t.cancel, 0, 1) {
		close(t.cancelChan)
		atomic.StoreInt32((*int32)(&t.state), int32(Cancelled))
	}
}

// ID returns the task's ID.
func (t *Task) ID() string {
	return t.id
}

// State returns the task's current state.
func (t *Task) State() TaskState {
	return TaskState(atomic.LoadInt32((*int32)(&t.state)))
}
