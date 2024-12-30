package timer

import (
	//"fmt"
	"runtime"
	"time"
)

// measurePerformance is a utility function that measures the time and memory usage of the given function.
func MeasurePerformance(fn func()) (time.Duration, uint64, uint64) {
	// Record the start time
	start := time.Now()

	// Record memory stats before the function execution
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// Execute the function
	fn()

	// Record memory stats after the function execution
	var memEnd runtime.MemStats
	runtime.ReadMemStats(&memEnd)

	// Calculate execution time
	duration := time.Since(start)

	// Calculate memory usage
	allocatedMemory := memEnd.Alloc - memStart.Alloc // Bytes allocated during the function execution
	totalMemory := memEnd.TotalAlloc // Total bytes allocated since the program started

	return duration, allocatedMemory, totalMemory
}

// sampleFunction is a sample function to measure
// func sampleFunction() {
// 	// Simulate some work (e.g., allocating memory)
// 	slice := make([]int, 1e6) // Allocate a slice of 1 million integers
// 	for i := 0; i < len(slice); i++ {
// 		slice[i] = i
// 	}
// 	// Here we could have more logic, but we'll keep it simple
// 	_ = slice // Prevent the compiler from optimizing away the slice
// }

