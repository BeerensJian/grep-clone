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

	// wait for workers in goroutine to still be able to print while working
	blockworkersWg := make(chan struct{})
	go func() {
		workersWg.Wait()
		close(blockworkersWg) // when channel is closed and empty, returns default value for its type
	}()

	var displayWg sync.WaitGroup

	displayWg.Add(1)
	go func() {
		for {
			select {
			case result := <-results:
				fmt.Printf("%v[%v]:%v\n", result.Path, result.LineNum, result.Line)
			case <-blockworkersWg:
				if len(results) == 0 {
					displayWg.Done()
					return
				}
			}
		}
	}()

	displayWg.Wait()

}
