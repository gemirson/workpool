package scheduler

import (
	"testing"
	"time"

	"github.com/gemirson/workpool/pkg/task"
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

	
}

func TestWorkpool_Cancel(t *testing.T) {
	wp := NewWorkpool(4)

	task1 := task.NewTask(func() {
		time.Sleep(4 * time.Second)
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

	
}

func TestWorkpool_Shutdown(t *testing.T) {
	wp := NewWorkpool(4)

	task1 := task.NewTask(func() {
		time.Sleep(10 * time.Second)
	})
	task2 := task.NewTask(func() {
		time.Sleep(10 * time.Second)
	})

	wp.Submit(task1, 1)
	wp.Submit(task2, 2)

	go func() {
		time.Sleep(1 * time.Second)
		wp.Shutdown()
	}()

	time.Sleep(3 * time.Second)

	assert.Equal(t, task.Running, task1.State())
	assert.Equal(t, task.Running, task2.State())
}
