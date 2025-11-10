package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// Fan-out/Fan-in Pattern Example: Distributed Computation
// Real-world use case: Calculating prime numbers across multiple goroutines
// and aggregating results

// isPrime checks if a number is prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// worker processes numbers and sends primes to results channel (Fan-out)
func primeWorker(id int, numbers <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for num := range numbers {
		if isPrime(num) {
			fmt.Printf("Worker %d: Found prime: %d\n", id, num)
			results <- num
		}
	}
}

// aggregator collects results from all workers (Fan-in)
func aggregator(results <-chan int, done chan<- []int) {
	primes := []int{}
	
	for prime := range results {
		primes = append(primes, prime)
	}
	
	done <- primes
}

func main() {
	fmt.Println("=== Fan-out/Fan-in Pattern Example ===")
	fmt.Println("Distributed Prime Number Calculation\n")
	
	start := time.Now()
	
	// Range of numbers to check
	startNum := 1
	endNum := 1000
	
	// Create channels
	numbers := make(chan int, 100)
	results := make(chan int, 100)
	done := make(chan []int)
	
	// Fan-out: Start multiple workers
	numWorkers := 5
	var wg sync.WaitGroup
	
	fmt.Printf("Starting %d workers to find primes between %d and %d...\n\n",
		numWorkers, startNum, endNum)
	
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go primeWorker(i, numbers, results, &wg)
	}
	
	// Fan-in: Start aggregator
	go aggregator(results, done)
	
	// Send numbers to workers
	go func() {
		for i := startNum; i <= endNum; i++ {
			numbers <- i
		}
		close(numbers)
	}()
	
	// Wait for workers to finish and close results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Get aggregated results
	primes := <-done
	
	elapsed := time.Since(start)
	
	fmt.Printf("\n=== Results ===\n")
	fmt.Printf("Found %d prime numbers in %v\n", len(primes), elapsed)
	fmt.Printf("First 10 primes: %v\n", primes[:int(math.Min(10, float64(len(primes))))])
	if len(primes) > 10 {
		fmt.Printf("Last 10 primes: %v\n", primes[len(primes)-10:])
	}
}
