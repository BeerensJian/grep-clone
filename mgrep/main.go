package main

import (
	"fmt"
	"mgrep/worker"
	"mgrep/worklist"
	"os"
	"path/filepath"
	"sync"
)

func GetAllFiles(wl *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			GetAllFiles(wl, nextPath)
		} else {
			wl.Add(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}

}

func main() {
	searchstring := os.Args[1]
	directory := os.Args[2]

	// search directory for all files and subdirectories
	// if subdirectory search subdirectory for all files (recursion?)

	wl := worklist.New(100)

	results := make(chan worker.Result, 100)

	numWorkers := 10

	var workersWg sync.WaitGroup

	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		GetAllFiles(&wl, directory)
		wl.Finalize(numWorkers)
	}()

	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for {
				work := wl.Next()
				if work.Path != "" {
					workerResult := worker.FindInFile(searchstring, work.Path)
					if workerResult != nil {
						for _, result := range workerResult.Inner {
							results <- result
						}
					}
				} else {
					return
				}
			}
		}()
	}

	// for each file search for substring in file and count the lines
	// if found put the found string into a channel and report it to the terminal

}
