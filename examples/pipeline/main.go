package main

import (
	"fmt"
	"strings"
)

// Pipeline Pattern Example: Data Transformation Pipeline
// Real-world use case: Processing log entries through multiple stages

type LogEntry struct {
	Raw       string
	Cleaned   string
	Parsed    map[string]string
	Enriched  map[string]string
}

// stage1: Read and clean log entries
func cleanLogs(logs <-chan string) <-chan LogEntry {
	out := make(chan LogEntry)
	
	go func() {
		defer close(out)
		for raw := range logs {
			entry := LogEntry{
				Raw:     raw,
				Cleaned: strings.TrimSpace(raw),
			}
			out <- entry
		}
	}()
	
	return out
}

// stage2: Parse log entries
func parseLogs(entries <-chan LogEntry) <-chan LogEntry {
	out := make(chan LogEntry)
	
	go func() {
		defer close(out)
		for entry := range entries {
			parts := strings.Split(entry.Cleaned, "|")
			parsed := make(map[string]string)
			
			if len(parts) >= 3 {
				parsed["timestamp"] = parts[0]
				parsed["level"] = parts[1]
				parsed["message"] = parts[2]
			}
			
			entry.Parsed = parsed
			out <- entry
		}
	}()
	
	return out
}

// stage3: Enrich log entries with metadata
func enrichLogs(entries <-chan LogEntry) <-chan LogEntry {
	out := make(chan LogEntry)
	
	go func() {
		defer close(out)
		for entry := range entries {
			enriched := make(map[string]string)
			for k, v := range entry.Parsed {
				enriched[k] = v
			}
			
			// Add enrichment
			if level, ok := entry.Parsed["level"]; ok {
				switch level {
				case "ERROR":
					enriched["severity"] = "high"
				case "WARN":
					enriched["severity"] = "medium"
				default:
					enriched["severity"] = "low"
				}
			}
			
			entry.Enriched = enriched
			out <- entry
		}
	}()
	
	return out
}

func main() {
	fmt.Println("=== Pipeline Pattern Example ===")
	fmt.Println("Log Processing Pipeline")
	fmt.Println()
	
	// Sample log entries
	rawLogs := []string{
		"2024-01-10 10:00:00|INFO|Application started successfully",
		"2024-01-10 10:01:23|ERROR|Database connection failed",
		"2024-01-10 10:02:45|WARN|High memory usage detected",
		"2024-01-10 10:03:12|INFO|User logged in",
		"2024-01-10 10:04:56|ERROR|API request timeout",
	}
	
	// Create input channel
	input := make(chan string)
	
	// Build the pipeline
	stage1 := cleanLogs(input)
	stage2 := parseLogs(stage1)
	stage3 := enrichLogs(stage2)
	
	// Feed data into the pipeline
	go func() {
		for _, log := range rawLogs {
			input <- log
		}
		close(input)
	}()
	
	// Consume pipeline output
	fmt.Println("Processing logs through pipeline...")
	fmt.Println()
	for entry := range stage3 {
		fmt.Printf("Raw: %s\n", entry.Raw)
		fmt.Printf("Parsed: %v\n", entry.Parsed)
		fmt.Printf("Enriched: %v\n\n", entry.Enriched)
	}
	
	fmt.Println("Pipeline processing completed!")
}
