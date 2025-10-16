package main

import (
	"os"
	"sync"
	"testing"
)

func TestFindWordInFile(t *testing.T) {
	content := []byte("hello world\nthis is a test\nfind the word hello here\nend of file")
	tmpfile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	var wg sync.WaitGroup
	ch := make(chan FindInfo, 1)
	targetWord := "hello"
	expectedLine := 2

	wg.Add(1)
	go FindWordInFile(targetWord, tmpfile.Name(), ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := <-ch

	if len(result.lines) != expectedLine {
		t.Errorf("Expected to find %d lines, but got %d", expectedLine, len(result.lines))
	}
	if result.lines[0].lineNo != 1 {
		t.Errorf("Expected line number to be 1, but got %d", result.lines[0].lineNo)
	}
	if result.lines[0].line != "hello world" {
		t.Errorf("Expected line content to be 'hello wrold', but got '%s'", result.lines[0].line)
	}
	if result.lines[1].lineNo != 3 {
		t.Errorf("Expected line content to be 3, but got '%d'", result.lines[1].lineNo)
	}
}