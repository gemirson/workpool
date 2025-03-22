package main

import (
	"fmt"
	"time"

	"github.com/gemirson/workpool/internal/scheduler"
	"github.com/gemirson/workpool/internal/task"
)

func main() {
	wp := scheduler.NewWorkpool(4) // Create a workpool with 4 workers

	task1 := task.NewTask(func() {
		fmt.Println("Task 1 executed")
	})
	task2 := task.NewTask(func() {
		fmt.Println("Task 2 executed")
	})

	wp.Submit(task1, 1) // Submit task1 with priority 1
	wp.Submit(task2, 2) // Submit task2 with priority 2

	time.Sleep(1 * time.Second)
	wp.Cancel(task1.ID()) // Cancel task1

	fmt.Printf("Task 1 state: %v\n", task1.State())
	fmt.Printf("Task 2 state: %v\n", task2.State())

	wp.Shutdown() // Shutdown the workpool
}
