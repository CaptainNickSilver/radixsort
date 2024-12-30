package radix_in_memory

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	
)

// radixSort sorts an array of integers using the Radix Sort algorithm.
func radixSort(arr []int) []int {
	max := 999999			// instead of using getMax which does a pass through the slice, we simply predetermine that the file has six digit numbers

	// Perform counting sort for each digit (1's, 10's, 100's, ..., 100000's)
	for exp := 1; max/exp > 0; exp *= 10 {
		countingSort(arr, exp)
	}
	return arr
}

// // getMax returns the maximum value in the array.
// func getMax(arr []int) int {
// 	max := arr[0]
// 	for _, num := range arr {
// 		if num > max {
// 			max = num
// 		}
// 	}
// 	return max
// }

// countingSort is a helper function that sorts the array based on the digit represented by exp.
func countingSort(arr []int, exp int) {
	n := len(arr)
	output := make([]int, n) // output array
	count := make([]int, 10) // count array for digits (0 to 9)

	// Store count of occurrences in count[]
	for i := 0; i < n; i++ {
		index := (arr[i] / exp) % 10
		count[index]++
	}

	// Change count[i] so that it contains the actual position of this digit in output[]
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	// Build the output array
	for i := n - 1; i >= 0; i-- {
		index := (arr[i] / exp) % 10
		output[count[index]-1] = arr[i]
		count[index]--
	}

	// Copy the output array to arr, so that arr now contains sorted numbers
	for i := 0; i < n; i++ {
		arr[i] = output[i]
	}
}

// readIntegers reads six-digit integers from a file and returns them as a slice.
func readIntegers(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var integers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid number %s: %v", line, err)
		}
		integers = append(integers, num)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return integers, nil
}

// writeIntegers writes a slice of integers to a file.
func writeIntegers(filename string, integers []int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, num := range integers {
		_, err := writer.WriteString(fmt.Sprintf("%d\n", num))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func Radix_in_Memory() {

	var inputFile, outputFile string 
	
	if len(os.Args) > 3 {
		// read the input file and output filename from the command line
		// if the user fails to provide both, error and quit
		inputFile = os.Args[2]
		outputFile = os.Args[3]
	} else {
		fmt.Println("Command line must contain the id of the algorithm and both the inputfile and output file names")
		return
	}

	// Read integers from the input file
	integers, err := readIntegers(inputFile)
	if err != nil {
		fmt.Println("Error reading integers:", err)
		return
	}

	// Sort the integers using Radix Sort
	sortedIntegers := radixSort(integers)

	// Write the sorted integers to the output file
	err = writeIntegers(outputFile, sortedIntegers)
	if err != nil {
		fmt.Println("Error writing integers:", err)
		return
	}

	fmt.Println("Sorting completed successfully. Sorted integers written to", outputFile)
}