package main

import (
	"testing"
	"time"
)

// Test basic goroutine execution
func TestBasicGoroutine(t *testing.T) {
	done := make(chan bool)
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		done <- true
	}()
	
	select {
	case <-done:
		// Success
	case <-time.After(200 * time.Millisecond):
		t.Error("Goroutine did not complete in time")
	}
}

// Test multiple goroutines
func TestMultipleGoroutines(t *testing.T) {
	count := 5
	results := make(chan int, count)
	
	for i := 0; i < count; i++ {
		go func(id int) {
			results <- id
		}(i)
	}
	
	received := make(map[int]bool)
	for i := 0; i < count; i++ {
		id := <-results
		received[id] = true
	}
	
	if len(received) != count {
		t.Errorf("Expected %d results, got %d", count, len(received))
	}
}

// Benchmark goroutine creation
func BenchmarkGoroutineCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan bool)
		go func() {
			done <- true
		}()
		<-done
	}
}
