package main

import (
	"sync"
	"testing"
)

// Test WaitGroup basic functionality
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	counter := 0
	mu := sync.Mutex{}
	
	numGoroutines := 10
	wg.Add(numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	if counter != numGoroutines {
		t.Errorf("Expected counter to be %d, got %d", numGoroutines, counter)
	}
}

// Test concurrent increment with WaitGroup
func TestConcurrentIncrement(t *testing.T) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sum := 0
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			mu.Lock()
			sum += val
			mu.Unlock()
		}(i)
	}
	
	wg.Wait()
	
	expected := 4950 // Sum of 0 to 99
	if sum != expected {
		t.Errorf("Expected sum to be %d, got %d", expected, sum)
	}
}

// Benchmark WaitGroup
func BenchmarkWaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(10)
		
		for j := 0; j < 10; j++ {
			go func() {
				defer wg.Done()
			}()
		}
		
		wg.Wait()
	}
}
