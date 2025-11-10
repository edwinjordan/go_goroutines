package main

import (
	"sync"
	"testing"
	"time"
)

// Test worker pool processes all jobs
func TestWorkerPool(t *testing.T) {
	numJobs := 20
	numWorkers := 3
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	var wg sync.WaitGroup
	
	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- job * 2
			}
		}()
	}
	
	// Send jobs
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)
	
	// Wait for workers
	wg.Wait()
	close(results)
	
	// Verify results
	count := 0
	for range results {
		count++
	}
	
	if count != numJobs {
		t.Errorf("Expected %d results, got %d", numJobs, count)
	}
}

// Test worker pool with different worker counts
func TestWorkerPoolScaling(t *testing.T) {
	testCases := []int{1, 2, 4, 8}
	
	for _, numWorkers := range testCases {
		t.Run(string(rune('0'+numWorkers))+" workers", func(t *testing.T) {
			numJobs := 100
			jobs := make(chan int, numJobs)
			results := make(chan int, numJobs)
			
			var wg sync.WaitGroup
			
			for i := 0; i < numWorkers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for job := range jobs {
						time.Sleep(time.Millisecond)
						results <- job
					}
				}()
			}
			
			for i := 0; i < numJobs; i++ {
				jobs <- i
			}
			close(jobs)
			
			wg.Wait()
			close(results)
			
			count := 0
			for range results {
				count++
			}
			
			if count != numJobs {
				t.Errorf("Expected %d results with %d workers, got %d", 
					numJobs, numWorkers, count)
			}
		})
	}
}

// Benchmark worker pool
func BenchmarkWorkerPool(b *testing.B) {
	numWorkers := 4
	
	for i := 0; i < b.N; i++ {
		jobs := make(chan int, 100)
		results := make(chan int, 100)
		
		var wg sync.WaitGroup
		
		for j := 0; j < numWorkers; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobs {
					results <- job * 2
				}
			}()
		}
		
		for j := 0; j < 100; j++ {
			jobs <- j
		}
		close(jobs)
		
		wg.Wait()
		close(results)
	}
}
