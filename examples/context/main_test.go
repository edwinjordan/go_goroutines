package main

import (
	"context"
	"testing"
	"time"
)

// Test context timeout
func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	select {
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have timed out")
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %v", ctx.Err())
		}
	}
}

// Test context cancellation
func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	done := make(chan bool)
	
	go func() {
		<-ctx.Done()
		done <- true
	}()
	
	// Cancel the context
	cancel()
	
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Error("Context cancellation not received in time")
	}
}

// Test context with value
func TestContextWithValue(t *testing.T) {
	type key string
	ctx := context.WithValue(context.Background(), key("user"), "john")
	
	value := ctx.Value(key("user"))
	if value != "john" {
		t.Errorf("Expected 'john', got %v", value)
	}
}

// Test goroutine respects context cancellation
func TestGoroutineCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	started := make(chan bool)
	finished := make(chan bool)
	
	go func() {
		started <- true
		select {
		case <-ctx.Done():
			finished <- true
		case <-time.After(1 * time.Second):
			finished <- false
		}
	}()
	
	<-started
	cancel()
	
	result := <-finished
	if !result {
		t.Error("Goroutine did not respect context cancellation")
	}
}

// Benchmark context creation
func BenchmarkContextCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_ = ctx
	}
}
