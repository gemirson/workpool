package notification

import "fmt"

// Notifier handles task completion notifications.
type Notifier struct {
	// Add fields if necessary
}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Notify sends a notification that the task with the given ID has completed.
func (n *Notifier) Notify(taskID string) {
	// Implement the notification logic
	fmt.Printf("Task %s completed\n", taskID)
}
