package main

import (
	"os"
	"sync"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindWordInFile(t *testing.T) {
	content := []byte("hello world\nthis is a test\nfind the word hello here\nend of file")
	tmpfile, err := os.CreateTemp("", "testfile_*.txt")
	require.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err, "Failed to create temp file")
	
	err = tmpfile.Close()
	require.NoError(t, err, "Failed to close temp file")

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

	assert.Equal(t, expectedLine, len(result.lines), "Expected to find %d lines, but got %d", expectedLine, len(result.lines))

	assert.Equal(t, 1, result.lines[0].lineNo, "Expected line number to be 1, but got %d", result.lines[0].lineNo)
	assert.Equal(t, "hello world", result.lines[0].line, "Expected line content to be 'hello wrold', but got '%s'", result.lines[0].line)
	
	assert.Equal(t, 3, result.lines[1].lineNo, "Second match line number should be 3")
}