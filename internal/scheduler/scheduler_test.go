package scheduler

import (
	"testing"
	"time"

	"github.com/gemirson/workpool/internal/task"
	"github.com/stretchr/testify/assert"
)

func TestWorkpool_Submit(t *testing.T) {
	wp := NewWorkpool(4)

	task1 := task.NewTask(func() {
		time.Sleep(1 * time.Second)
	})
	task2 := task.NewTask(func() {
		time.Sleep(1 * time.Second)
	})

	wp.Submit(task1, 1)
	wp.Submit(task2, 2)

	time.Sleep(2 * time.Second)

	assert.Equal(t, task.Completed, task1.State())
	assert.Equal(t, task.Completed, task2.State())

	wp.Shutdown()
}

func TestWorkpool_Cancel(t *testing.T) {
	wp := NewWorkpool(4)

	task1 := task.NewTask(func() {
		time.Sleep(1 * time.Second)
	})
	task2 := task.NewTask(func() {
		time.Sleep(1 * time.Second)
	})

	wp.Submit(task1, 1)
	wp.Submit(task2, 2)

	wp.Cancel(task1.ID())

	time.Sleep(2 * time.Second)

	assert.Equal(t, task.Cancelled, task1.State())
	assert.Equal(t, task.Completed, task2.State())

	wp.Shutdown()
}

func TestWorkpool_Shutdown(t *testing.T) {
	wp := NewWorkpool(4)

	task1 := task.NewTask(func() {
		time.Sleep(1 * time.Second)
	})
	task2 := task.NewTask(func() {
		time.Sleep(1 * time.Second)
	})

	wp.Submit(task1, 1)
	wp.Submit(task2, 2)

	wp.Shutdown()

	assert.Equal(t, task.Completed, task1.State())
	assert.Equal(t, task.Completed, task2.State())
	https://github.com/gemirson/workpool
	https://github.com/gemirson/workpool
	https://github.com/gemirson/workpool
}
