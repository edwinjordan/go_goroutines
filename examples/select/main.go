package main

import (
	"fmt"
	"time"
)

// Select Statement Example: Timeout and Multiplexing
// Real-world use case: Handling multiple channels with timeouts

// fetchData simulates a data fetch operation
func fetchData(source string, delay time.Duration) <-chan string {
	ch := make(chan string)
	
	go func() {
		time.Sleep(delay)
		ch <- fmt.Sprintf("Data from %s", source)
	}()
	
	return ch
}

// monitorSystem simulates system monitoring
func monitorSystem() <-chan string {
	ch := make(chan string)
	
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(time.Second)
			ch <- fmt.Sprintf("System check #%d: OK", i)
		}
		close(ch)
	}()
	
	return ch
}

func main() {
	fmt.Println("=== Select Statement Example ===")
	fmt.Println("Timeout Handling and Channel Multiplexing")
	fmt.Println()
	
	// Example 1: Timeout handling
	fmt.Println("Example 1: Fetch with timeout")
	
	data1 := fetchData("API-1", 2*time.Second)
	data2 := fetchData("API-2", 5*time.Second) // This will timeout
	
	select {
	case result := <-data1:
		fmt.Printf("✓ Received: %s\n", result)
	case <-time.After(3 * time.Second):
		fmt.Println("✗ Timeout waiting for API-1")
	}
	
	select {
	case result := <-data2:
		fmt.Printf("✓ Received: %s\n", result)
	case <-time.After(3 * time.Second):
		fmt.Println("✗ Timeout waiting for API-2 (expected)")
		fmt.Println()
	}
	
	// Example 2: Multiplexing multiple channels
	fmt.Println("Example 2: Multiplexing channels")
	
	channel1 := make(chan string)
	channel2 := make(chan string)
	quit := make(chan bool)
	
	// Send data to channels at different intervals
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(time.Millisecond * 500)
			channel1 <- fmt.Sprintf("Channel 1: Message %d", i)
		}
	}()
	
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(time.Millisecond * 700)
			channel2 <- fmt.Sprintf("Channel 2: Message %d", i)
		}
	}()
	
	go func() {
		time.Sleep(3 * time.Second)
		quit <- true
	}()
	
	// Multiplex the channels
	messageCount := 0
	for {
		select {
		case msg := <-channel1:
			fmt.Printf("  %s\n", msg)
			messageCount++
		case msg := <-channel2:
			fmt.Printf("  %s\n", msg)
			messageCount++
		case <-quit:
			fmt.Printf("\nReceived quit signal. Processed %d messages\n", messageCount)
			return
		case <-time.After(4 * time.Second):
			fmt.Println("\nNo activity for 4 seconds, exiting...")
			return
		}
	}
}
