package main

import (
	"os"
)

func main() {
	searchstring := os.Args[1]
	directory := os.Args[2]

	// search directory for all files and subdirectories
	// if subdirectory search subdirectory for all files (recursion?)
	// for each file search for substring in file and count the lines
	// if found put the found string into a channel and report it to the terminal

}
