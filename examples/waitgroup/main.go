package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// WaitGroup Example: Concurrent API Calls
// Real-world use case: Fetching data from multiple APIs concurrently

type APIResult struct {
	URL        string
	StatusCode int
	Length     int
	Duration   time.Duration
	Error      error
}

// fetchURL makes an HTTP request and returns the result
func fetchURL(url string, wg *sync.WaitGroup, results chan<- APIResult) {
	defer wg.Done()
	
	start := time.Now()
	result := APIResult{URL: url}
	
	fmt.Printf("Fetching: %s\n", url)
	
	resp, err := http.Get(url)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		results <- result
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		results <- result
		return
	}
	
	result.StatusCode = resp.StatusCode
	result.Length = len(body)
	result.Duration = time.Since(start)
	results <- result
}

func main() {
	fmt.Println("=== WaitGroup Example ===")
	fmt.Println("Concurrent API Calls")
	fmt.Println()
	
	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/users/1",
		"https://jsonplaceholder.typicode.com/comments/1",
		"https://jsonplaceholder.typicode.com/albums/1",
		"https://jsonplaceholder.typicode.com/photos/1",
	}
	
	var wg sync.WaitGroup
	results := make(chan APIResult, len(urls))
	
	// Launch goroutines for each URL
	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg, results)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	close(results)
	
	// Display results
	fmt.Println("\n=== Results ===")
	for result := range results {
		if result.Error != nil {
			fmt.Printf("❌ %s - Error: %v\n", result.URL, result.Error)
		} else {
			fmt.Printf("✓ %s - Status: %d, Size: %d bytes, Time: %v\n",
				result.URL, result.StatusCode, result.Length, result.Duration)
		}
	}
}
