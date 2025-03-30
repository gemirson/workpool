package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask_Execute(t *testing.T) {
	executed := false
	tk := NewTask(func() {
		executed = true
	})

	tk.Execute()

	assert.True(t, executed)
	assert.Equal(t, Completed, tk.State())
}

func TestTask_Cancel(t *testing.T) {
	executed := false
	tk := NewTask(func() {
		executed = true
		time.Sleep(1 * time.Second)
	})

	tk.Cancel()
	tk.Execute()

	assert.False(t, executed)
	assert.Equal(t, Cancelled, tk.State())
}
