package main

import (
	"fmt"

	"github.com/gemirson/workpool/pkg/scheduler"
	"github.com/gemirson/workpool/pkg/task"
)

// This example demonstrates how to create a workpool, submit a task, and shut it down.
func main() {
	wp := scheduler.NewWorkpool(10) // Create a workpool with 10 workers

	task1 := task.NewTask(func() {
		fmt.Println("Task 1 executed")
	})

	wp.Submit(task1) // Submit a task to the workpool
	wp.Shutdown()    // Shutdown the workpool
}
