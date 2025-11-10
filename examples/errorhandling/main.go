package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Error Handling Example: Concurrent Error Aggregation
// Real-world use case: Running multiple tasks and collecting all errors

// Task represents a unit of work that may fail
type Task struct {
	ID   int
	Name string
}

// TaskResult holds the result or error of a task
type TaskResult struct {
	TaskID int
	Name   string
	Error  error
	Data   string
}

// processTask simulates processing a task that may fail
func processTask(task Task) TaskResult {
	// Simulate processing time
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	
	result := TaskResult{
		TaskID: task.ID,
		Name:   task.Name,
	}
	
	// Simulate random failures (30% failure rate)
	if rand.Float32() < 0.3 {
		result.Error = errors.New("processing failed due to random error")
		return result
	}
	
	result.Data = fmt.Sprintf("Processed data for %s", task.Name)
	return result
}

// worker processes tasks and sends results
func worker(id int, tasks <-chan Task, results chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for task := range tasks {
		fmt.Printf("Worker %d: Processing task %d (%s)\n", id, task.ID, task.Name)
		result := processTask(task)
		results <- result
	}
}

// ErrorCollector collects and manages errors from concurrent operations
type ErrorCollector struct {
	mu     sync.Mutex
	errors []error
}

// Add adds an error to the collector
func (ec *ErrorCollector) Add(err error) {
	if err != nil {
		ec.mu.Lock()
		ec.errors = append(ec.errors, err)
		ec.mu.Unlock()
	}
}

// HasErrors returns true if there are any errors
func (ec *ErrorCollector) HasErrors() bool {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return len(ec.errors) > 0
}

// GetErrors returns all collected errors
func (ec *ErrorCollector) GetErrors() []error {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	return append([]error{}, ec.errors...)
}

func main() {
	fmt.Println("=== Error Handling Example ===")
	fmt.Println("Concurrent Error Aggregation\n")
	
	rand.Seed(time.Now().UnixNano())
	
	// Create tasks
	tasks := []Task{
		{ID: 1, Name: "ValidateData"},
		{ID: 2, Name: "TransformData"},
		{ID: 3, Name: "EnrichData"},
		{ID: 4, Name: "FilterData"},
		{ID: 5, Name: "AggregateData"},
		{ID: 6, Name: "ExportData"},
		{ID: 7, Name: "BackupData"},
		{ID: 8, Name: "ArchiveData"},
		{ID: 9, Name: "CleanupData"},
		{ID: 10, Name: "IndexData"},
	}
	
	// Create channels
	taskChan := make(chan Task, len(tasks))
	resultChan := make(chan TaskResult, len(tasks))
	
	// Start workers
	numWorkers := 3
	var wg sync.WaitGroup
	
	fmt.Printf("Starting %d workers...\n\n", numWorkers)
	
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChan, resultChan, &wg)
	}
	
	// Send tasks
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)
	
	// Wait for workers to finish
	wg.Wait()
	close(resultChan)
	
	// Collect results and errors
	errorCollector := &ErrorCollector{}
	successCount := 0
	failedTasks := []string{}
	
	fmt.Println("\n=== Results ===")
	for result := range resultChan {
		if result.Error != nil {
			fmt.Printf("❌ Task %d (%s): %v\n", result.TaskID, result.Name, result.Error)
			errorCollector.Add(fmt.Errorf("task %d (%s): %w", result.TaskID, result.Name, result.Error))
			failedTasks = append(failedTasks, result.Name)
		} else {
			fmt.Printf("✓ Task %d (%s): Success - %s\n", result.TaskID, result.Name, result.Data)
			successCount++
		}
	}
	
	// Summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total tasks: %d\n", len(tasks))
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", len(failedTasks))
	
	if errorCollector.HasErrors() {
		fmt.Printf("\nFailed tasks: %v\n", failedTasks)
		fmt.Printf("Total errors collected: %d\n", len(errorCollector.GetErrors()))
	} else {
		fmt.Println("\n✓ All tasks completed successfully!")
	}
}
