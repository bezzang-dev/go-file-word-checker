package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LineInfo holds the line number and content where a word is find.
type LineInfo struct {
	lineNo int
	line   string
}

// FindInfo holds the filename and all the lines that contain the target word.
type FindInfo struct {
	filename string
	lines    []LineInfo
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./go-file-word-checker <word> <file/path pattern>")
		fmt.Println("Example: ./go-file-word-checker 'hello' '*.txt'")
		return
	}

	word := os.Args[1]
	paths := os.Args[2:]

	findInfos := FindWordInAllFiles(word, paths)
	
	for _, findInfo := range findInfos {
		fmt.Println(findInfo.filename)
		fmt.Println("----------------------")
		for _, lineInfo := range findInfo.lines {
			fmt.Println("\t", lineInfo.lineNo, "\t", lineInfo.line)
		}
		fmt.Println("--------------------------")
		fmt.Println()
	}
}

// FindWordInAllFiles searches for a word across multiple file path patterns concurrently.
func FindWordInAllFiles(word string, paths []string) []FindInfo {
	var wg sync.WaitGroup
	findInfos := []FindInfo{}
	ch := make(chan FindInfo)

	for _, path := range paths {

		// Get all files that match the pattern
		filelist, err := filepath.Glob(path)
		if err != nil {
			fmt.Printf("Error with path pattern %s: %v\n", path, err)
			continue
		}

		for _, filename := range filelist {
			wg.Add(1)
			go FindWordInFile(word, filename, ch, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for findInfo := range ch {
		findInfos = append(findInfos, findInfo)
	}

	return findInfos
}

func FindWordInFile(word, filename string, ch chan<- FindInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	findInfo := FindInfo{filename, []LineInfo{}}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Cannot open file %s: %v\n", filename, err)
		ch <- findInfo
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := scanner.Text()
		if strings.Contains(line, word) {
			findInfo.lines = append(findInfo.lines, LineInfo{lineNo, line})
		}
	}
	ch <- findInfo
}