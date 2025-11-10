package main

import (
	"fmt"
	"sync"
	"time"
)

// Rate Limiting Example: API Rate Limiter
// Real-world use case: Limiting API requests to avoid overwhelming a service

// RateLimiter controls the rate of goroutine execution
type RateLimiter struct {
	ticker *time.Ticker
	tokens chan struct{}
}

// NewRateLimiter creates a rate limiter that allows 'rate' operations per second
func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
		tokens: make(chan struct{}, rate),
	}
	
	// Fill initial tokens
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}
	
	// Refill tokens
	go func() {
		for range rl.ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Token bucket is full
			}
		}
	}()
	
	return rl
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
	<-rl.tokens
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.tokens)
}

// apiRequest simulates an API request
func apiRequest(id int, url string, limiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Wait for rate limiter
	limiter.Wait()
	
	start := time.Now()
	fmt.Printf("[%s] Request %d: Calling %s\n", start.Format("15:04:05"), id, url)
	
	// Simulate API call
	time.Sleep(time.Millisecond * 100)
	
	elapsed := time.Since(start)
	fmt.Printf("[%s] Request %d: Completed in %v\n", time.Now().Format("15:04:05"), id, elapsed)
}

func main() {
	fmt.Println("=== Rate Limiting Example ===")
	fmt.Println("API Rate Limiter\n")
	
	// Create rate limiter: 5 requests per second
	requestsPerSecond := 5
	limiter := NewRateLimiter(requestsPerSecond)
	defer limiter.Stop()
	
	fmt.Printf("Rate limit: %d requests per second\n", requestsPerSecond)
	fmt.Println("Sending 15 requests...\n")
	
	// Simulate 15 API requests
	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/posts",
		"https://api.example.com/comments",
		"https://api.example.com/albums",
		"https://api.example.com/photos",
		"https://api.example.com/todos",
		"https://api.example.com/users/1",
		"https://api.example.com/posts/1",
		"https://api.example.com/comments/1",
		"https://api.example.com/albums/1",
		"https://api.example.com/photos/1",
		"https://api.example.com/todos/1",
		"https://api.example.com/users/2",
		"https://api.example.com/posts/2",
		"https://api.example.com/comments/2",
	}
	
	var wg sync.WaitGroup
	start := time.Now()
	
	for i, url := range urls {
		wg.Add(1)
		go apiRequest(i+1, url, limiter, &wg)
	}
	
	wg.Wait()
	
	totalTime := time.Since(start)
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total requests: %d\n", len(urls))
	fmt.Printf("Total time: %v\n", totalTime)
	fmt.Printf("Average time per request: %v\n", totalTime/time.Duration(len(urls)))
	fmt.Printf("Expected minimum time: ~%v (due to rate limiting)\n", 
		time.Duration(len(urls)/requestsPerSecond)*time.Second)
}
