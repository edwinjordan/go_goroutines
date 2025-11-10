package main

import (
	"testing"
)

// Test channel communication
func TestChannelCommunication(t *testing.T) {
	ch := make(chan int)
	
	go func() {
		ch <- 42
	}()
	
	result := <-ch
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

// Test buffered channel
func TestBufferedChannel(t *testing.T) {
	ch := make(chan int, 3)
	
	// Should not block
	ch <- 1
	ch <- 2
	ch <- 3
	
	if len(ch) != 3 {
		t.Errorf("Expected buffer length 3, got %d", len(ch))
	}
	
	// Receive all
	for i := 1; i <= 3; i++ {
		result := <-ch
		if result != i {
			t.Errorf("Expected %d, got %d", i, result)
		}
	}
}

// Test channel close
func TestChannelClose(t *testing.T) {
	ch := make(chan int)
	
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()
	
	count := 0
	for range ch {
		count++
	}
	
	if count != 3 {
		t.Errorf("Expected to receive 3 values, got %d", count)
	}
}

// Benchmark channel operations
func BenchmarkChannelSendReceive(b *testing.B) {
	ch := make(chan int)
	
	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
	}()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
}
