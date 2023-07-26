package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line    string
	LineNum int
	Path    string
}

type Results struct {
	Inner []Result
}

func NewResult(line string, lineNum int, path string) Result {
	return Result{line, lineNum, path}
}

func FindInFile(searchTerm string, path string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}

	results := make([]Result, 0)

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber += lineNumber
		line := scanner.Text()
		if strings.Contains(line, searchTerm) {
			results = append(results, Result{Line: line, LineNum: lineNumber, Path: path})
		}

	}

	return &Results{results}
}
