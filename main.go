package main

import (
	"fmt"
	"os"
	"time"

	rim "github.com/captainnicksilver/radixsort/radix_in_memory"
	rod "github.com/captainnicksilver/radixsort/radix_on_disk"
	"github.com/captainnicksilver/radixsort/timer"
)

func main() {


	var duration time.Duration
	var allocatedMemory, totalMemory uint64
	var inputFile, outputFile string

	// Measure the performance of the indicated function
	if len(os.Args) > 3 {
		id := os.Args[1] // id of the algorithm to test
		inputFile = os.Args[2]
		outputFile = os.Args[3]
		switch id {
		case "RIM":
			duration, allocatedMemory, totalMemory = timer.MeasurePerformance(rim.Radix_in_Memory, inputFile, outputFile)
		case "ROD":
			duration, allocatedMemory, totalMemory = timer.MeasurePerformance( rod.Radix_On_Disk, inputFile, outputFile)
		default:
			fmt.Println("The id provided does not match any of the coded algorithms.  Use RIM for Radix In Memory or ROD for Radix on Disk")
			return
		}

	} else {
		fmt.Println("The input line must contain the id of the algorithm and both the input and output files")
		return
	}

	

	// Print the results
	fmt.Printf("Execution time: %v\n", duration)
	fmt.Printf("Allocated memory during execution: %d bytes\n", allocatedMemory)
	fmt.Printf("Total memory allocated since start: %d bytes\n", totalMemory)
}