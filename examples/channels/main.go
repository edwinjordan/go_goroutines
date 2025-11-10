package main

import (
	"fmt"
	"time"
)

// Goroutines with Channels Example: Producer-Consumer Pattern
// Real-world use case: Message queue system where producers send messages 
// and consumers process them

// producer generates data and sends it to a channel
func producer(id int, jobs chan<- int) {
	for i := 1; i <= 5; i++ {
		job := id*10 + i
		fmt.Printf("Producer %d: Sending job %d\n", id, job)
		jobs <- job
		time.Sleep(time.Millisecond * 500)
	}
}

// consumer receives data from a channel and processes it
func consumer(id int, jobs <-chan int, done chan<- bool) {
	for job := range jobs {
		fmt.Printf("  Consumer %d: Processing job %d\n", id, job)
		time.Sleep(time.Millisecond * 300) // Simulate work
		fmt.Printf("  Consumer %d: Completed job %d\n", id, job)
	}
	done <- true
}

func main() {
	fmt.Println("=== Goroutines with Channels Example ===")
	fmt.Println("Producer-Consumer Pattern")
	fmt.Println()
	
	// Create channels
	jobs := make(chan int, 10) // Buffered channel for jobs
	done := make(chan bool)     // Channel to signal completion
	
	// Start 2 producers
	go producer(1, jobs)
	go producer(2, jobs)
	
	// Start 3 consumers
	numConsumers := 3
	for i := 1; i <= numConsumers; i++ {
		go consumer(i, jobs, done)
	}
	
	// Wait for producers to finish
	time.Sleep(3 * time.Second)
	close(jobs) // Close jobs channel when no more jobs
	
	// Wait for all consumers to finish
	for i := 0; i < numConsumers; i++ {
		<-done
	}
	
	fmt.Println("\nAll jobs processed!")
}
