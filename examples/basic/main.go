package main

import (
	"fmt"
	"time"
)

// Basic Goroutines Example: Concurrent File Processing
// Real-world use case: Processing multiple files concurrently to count lines

// processFile simulates processing a file (counting lines)
func processFile(filename string) {
	fmt.Printf("Starting to process: %s\n", filename)
	
	// Simulate file processing time
	time.Sleep(time.Duration(len(filename)%3+1) * time.Second)
	
	// Simulate line counting
	lineCount := len(filename) * 10 // Mock line count
	
	fmt.Printf("Finished processing %s: %d lines\n", filename, lineCount)
}

func main() {
	files := []string{
		"data.csv",
		"report.txt",
		"config.json",
		"log.txt",
		"users.db",
	}
	
	fmt.Println("=== Basic Goroutines Example ===")
	fmt.Println("Processing files concurrently...")
	
	// Launch goroutines for concurrent file processing
	for _, file := range files {
		go processFile(file)
	}
	
	// Wait for all goroutines to finish
	// Note: This is a simple example - in production use WaitGroup or channels
	time.Sleep(5 * time.Second)
	
	fmt.Println("\nAll files processed!")
}
