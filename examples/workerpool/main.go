package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker Pool Pattern Example: Image Processing Simulation
// Real-world use case: Processing images concurrently with a fixed number of workers

type Image struct {
	ID       int
	Name     string
	Size     int
	Processed bool
}

// worker processes images from the jobs channel
func worker(id int, jobs <-chan Image, results chan<- Image, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for image := range jobs {
		fmt.Printf("Worker %d: Processing image %s (ID: %d, Size: %d KB)\n",
			id, image.Name, image.ID, image.Size)
		
		// Simulate image processing time (proportional to size)
		processingTime := time.Duration(image.Size) * time.Millisecond * 10
		time.Sleep(processingTime)
		
		// Mark as processed
		image.Processed = true
		
		fmt.Printf("Worker %d: Completed image %s\n", id, image.Name)
		results <- image
	}
}

func main() {
	fmt.Println("=== Worker Pool Pattern Example ===")
	fmt.Println("Image Processing Simulation")
	fmt.Println()
	
	// Generate sample images
	images := []Image{
		{ID: 1, Name: "photo1.jpg", Size: 15},
		{ID: 2, Name: "photo2.jpg", Size: 20},
		{ID: 3, Name: "photo3.jpg", Size: 10},
		{ID: 4, Name: "photo4.jpg", Size: 25},
		{ID: 5, Name: "photo5.jpg", Size: 12},
		{ID: 6, Name: "photo6.jpg", Size: 18},
		{ID: 7, Name: "photo7.jpg", Size: 22},
		{ID: 8, Name: "photo8.jpg", Size: 14},
	}
	
	// Create channels
	jobs := make(chan Image, len(images))
	results := make(chan Image, len(images))
	
	// Start worker pool with 3 workers
	numWorkers := 3
	var wg sync.WaitGroup
	
	fmt.Printf("Starting %d workers...\n\n", numWorkers)
	
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}
	
	// Send jobs to workers
	for _, img := range images {
		jobs <- img
	}
	close(jobs)
	
	// Wait for all workers to finish
	wg.Wait()
	close(results)
	
	// Collect results
	fmt.Println("\n=== Processing Summary ===")
	processedCount := 0
	for img := range results {
		if img.Processed {
			processedCount++
		}
	}
	
	fmt.Printf("Total images: %d\n", len(images))
	fmt.Printf("Successfully processed: %d\n", processedCount)
}
