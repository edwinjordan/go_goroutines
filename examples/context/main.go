package main

import (
	"context"
	"fmt"
	"time"
)

// Context-based Cancellation Example: Graceful Shutdown
// Real-world use case: Long-running workers that can be cancelled gracefully

// longRunningTask simulates a long-running operation that respects context
func longRunningTask(ctx context.Context, id int, results chan<- string) {
	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			// Context cancelled
			results <- fmt.Sprintf("Worker %d: Cancelled at step %d (reason: %v)", id, i, ctx.Err())
			return
		case <-time.After(500 * time.Millisecond):
			// Continue working
			msg := fmt.Sprintf("Worker %d: Completed step %d/10", id, i)
			fmt.Println(msg)
		}
	}
	results <- fmt.Sprintf("Worker %d: Completed all steps successfully", id)
}

// dataProcessor simulates processing data with context
func dataProcessor(ctx context.Context, name string) {
	fmt.Printf("%s: Started\n", name)
	
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: Shutting down gracefully... (%v)\n", name, ctx.Err())
			// Cleanup operations here
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("%s: Cleanup completed\n", name)
			return
		case t := <-ticker.C:
			fmt.Printf("%s: Processing data at %s\n", name, t.Format("15:04:05"))
		}
	}
}

func main() {
	fmt.Println("=== Context-based Cancellation Example ===")
	fmt.Println("Graceful Shutdown of Goroutines\n")
	
	// Example 1: Timeout context
	fmt.Println("Example 1: Timeout Context")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()
	
	results := make(chan string, 2)
	
	go longRunningTask(ctx1, 1, results)
	go longRunningTask(ctx1, 2, results)
	
	// Wait for results
	for i := 0; i < 2; i++ {
		result := <-results
		fmt.Printf("Result: %s\n", result)
	}
	
	fmt.Println("\nExample 2: Manual Cancellation")
	
	// Example 2: Manual cancellation
	ctx2, cancel2 := context.WithCancel(context.Background())
	
	// Start multiple processors
	go dataProcessor(ctx2, "Processor-A")
	go dataProcessor(ctx2, "Processor-B")
	go dataProcessor(ctx2, "Processor-C")
	
	// Let them run for a while
	time.Sleep(2500 * time.Millisecond)
	
	// Cancel the context
	fmt.Println("\n>>> Initiating shutdown...")
	cancel2()
	
	// Give time for graceful shutdown
	time.Sleep(1 * time.Second)
	
	fmt.Println("\nAll goroutines stopped gracefully!")
}
