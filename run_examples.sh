#!/bin/bash

# Script to run all goroutines examples

echo "========================================="
echo "  Go Goroutines Examples Runner"
echo "========================================="
echo ""

# Array of examples
examples=(
    "basic"
    "channels"
    "waitgroup"
    "workerpool"
    "pipeline"
    "fanout"
    "select"
    "context"
    "ratelimit"
    "errorhandling"
)

# Run each example
for example in "${examples[@]}"; do
    echo "----------------------------------------"
    echo "Running: $example"
    echo "----------------------------------------"
    timeout 15 go run examples/$example/main.go
    echo ""
    echo "Press Enter to continue to next example..."
    read -t 2
    echo ""
done

echo "========================================="
echo "  All examples completed!"
echo "========================================="
