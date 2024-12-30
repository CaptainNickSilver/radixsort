package radix_on_disk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// ---------------------------
// setup_output_files - create the list of output files for this pass
// ---------------------------
func setup_output_files(  outPath string, pass int) []string {

	outlist := make([]string,10)
	
	for ct := 0; ct < 10; ct++ {
		outlist[ct] = filepath.Join(outPath, "pass"+strconv.Itoa(pass)+"_"+strconv.Itoa(ct)+".txt")
	}
	return outlist
}

// ---------------------------
// open_output_files - take a list of filenames and open each one, returning an array of file handles
// ---------------------------
func open_output_files(outflist []string) []*os.File {
	fHandles := make([]*os.File, 0, 10)
	for _, filename := range outflist {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error opening file(%s): %s\n", filename, err.Error())
			// Close any files that were opened successfully before returning the error.
			for _, f := range fHandles {
				f.Close()
			}
			return nil
		}
		fHandles = append(fHandles, file)
	}
	return fHandles
}

// ---------------------------
// close_output_files - close all of the file handles in an array
// ---------------------------

func close_output_files(outHandles []*os.File) {
	for _, fHandle := range outHandles {
		fHandle.Close()
	}
}

// ---------------------------
// process_num -- the heart of radix sort.  Select the correct file to push the number to based on Most Significant Digit
// ---------------------------
func process_num(passid int, valtosort int, outhandles []*os.File) {
	// passid is assumed to be between 0 and 5 inclusive.  All of the numbers are six digit integers
	var strrep string = strconv.Itoa(valtosort)
	byteArray := []byte(strrep) // convert to a byte array
	idx := int(byteArray[passid] - byte('0'))
	if len(outhandles) < idx {
		fmt.Printf("FAILURE to properly open files.  Abort operations")
		os.Exit(1)
	}
	outhandles[idx].Write(byteArray)
	outhandles[idx].WriteString("\r\n")
}

// ---------------------------
// Radix_On_Disk -- the main routine for this module -- single threaded but nothing is called into memory
// ---------------------------

func Radix_On_Disk() {
	var inputFile, outputFile, outPath string
	var rowcount int64 = 0
	infilelist := make([]string, 1, 10) // slice on array of input files, starting with length of 1
	var outfilelist []string
	var passcount int

	if len(os.Args) > 3 {
		// read the input file and output filename from the command line
		// if the user fails to provide both, error and quit
		infilelist[0] = os.Args[2]
		outputFile = os.Args[3]
		outPath = filepath.Dir(outputFile)
	} else {
		fmt.Println("Command line must contain the id of the algorithm, the inputfile and output file names")
		return
	}

	// input and output is a little weird on the disk version.  We will read from one file
	// on the first pass but we will read from 10 files on subsequent passes.  Since these
	// are six digit numbers, it will take a total of six passes.  Note that this is MSD Radix

	// to accomplish this, we will use an array of filenames.  The first pass, the array has a single
	// filename in it.  However, the write routine will create an array of filenames being written. That
	// new array becomes the input to the next round.  Note that, other than pass one, the files
	// being read are deleted after we hit EOF to save on disk space.  On the last pass, we will output to
	// a single file: the name of the file passed as a parameter

	for passcount = 0; passcount < 6; passcount++ {
		outfilelist = setup_output_files(outPath, passcount)   // intentionally stomp any preexisting allocation of output files
		outFileHandle := open_output_files(outfilelist)
		for _, inFile := range infilelist {
			//open the input file

			fileH, err := os.Open(inFile)
			if err != nil {
				fmt.Printf("Error opening input file for pass %d: %s error %s\n", passcount, inFile, err.Error())
				return
			}
			// do not defer the close, since we want to close and delete within the same loop

			// read and process a row into one of the outfiles

			scanner := bufio.NewScanner(fileH)
			for scanner.Scan() {
				line := scanner.Text()
				numVal, err := strconv.Atoi(line)
				if err != nil {
					fmt.Printf("invalid number %s: %v\n", line, err)
					return
				}
				if passcount == 0 {
					rowcount++
				} // count the rows in the input but only count once (on the first pass)
				process_num(passcount, numVal, outFileHandle)

			}
			// we have used up the input file.  close it and delete it
			fileH.Close()
			//os.Remove(inFile)		for the sake of debugging, keep all files
		}
		// after the pass is done, close all the outfiles,
		close_output_files(outFileHandle)

		// move the outfile slice to the infile array in preparation for next pass
		// note, this is a shallow copy (pointing to the same memory. On purpose.)
		infilelist = outfilelist
	}

	// now read all of the files in the infile list (that were just output by the last pass)
	// and copy the results to a single output file (concatenation)

	// set up the output file for io.copy
	outfile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error opening the output file:", err.Error())
		return
	}
	defer outfile.Close()

	for _, inputFile = range infilelist {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("error opening final pass input file (%s): %s\n", inputFile, err.Error())
			continue
		}
		// do not defer the close... we want to control the close and delete the file

		_, err = io.Copy(outfile, f)
		f.Close()
		if err == nil {  // if concat was successful, delete the source file
			os.Remove(inputFile)
		} else {
			fmt.Printf("error concatenating input file(%s): %s ", inputFile, err.Error())
		}

	}
	fmt.Printf("Sorted %d lines into %s\n", rowcount, outputFile)

}
