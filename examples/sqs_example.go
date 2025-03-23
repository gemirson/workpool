package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gemirson/workpool/pkg/scheduler"
	"github.com/gemirson/workpool/pkg/task"
)

const (
	region      = "us-west-2"
	queueURL    = "https://sqs.us-west-2.amazonaws.com/123456789012/MyQueue"
	workerCount = 4 // Set worker count to 4 for 4 available threads
)

func main() {
	// Initialize AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	// Create SQS service client
	svc := sqs.New(sess)

	// Create a workpool with 4 workers
	wp := scheduler.NewWorkpool(workerCount)

	// Function to process messages
	processMessage := func(msg *sqs.Message) {
		fmt.Printf("Processing message: %s\n", *msg.Body)
		// Simulate processing time
		time.Sleep(2 * time.Second)
	}

	// Function to receive messages from SQS
	receiveMessages := func() {
		for {
			result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: aws.Int64(10),
				WaitTimeSeconds:     aws.Int64(20),
			})
			if err != nil {
				log.Printf("Error receiving messages: %v", err)
				continue
			}

			for _, msg := range result.Messages {
				task := task.NewTask(func() {
					processMessage(msg)
				})
				wp.Submit(task, 1) // Submit task with priority 1
			}
		}
	}

	// Start receiving messages
	go receiveMessages()

	// Wait for user to terminate the program
	fmt.Println("Press Ctrl+C to exit...")
	select {}
}
