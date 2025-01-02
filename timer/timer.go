package timer

import (
	//"fmt"
	"runtime"
	"time"
)

// MeasurePerformance is a utility function that measures the time and memory usage of the given function.
// the only thing truly interesting about it is that we pass in a function and it's parameters seperately 
// so that we can call it.  This could be generalized further if desired.

func MeasurePerformance(fn func(string, string), infile string, outfile string) (time.Duration, uint64, uint64) {
	// Record the start time
	start := time.Now()

	// Record memory stats before the function execution
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// Execute the function
	fn(infile, outfile)

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

