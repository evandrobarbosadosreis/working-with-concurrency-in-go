package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printSomething(t *testing.T) {
	// save the standard output
	stdout := os.Stdout
	// creates a connected writer -> reader
	r, w, _ := os.Pipe()
	// defines my writter as the new standard output (so I can read it later)
	os.Stdout = w
	// creates the waitgroup
	var wg sync.WaitGroup
	// setup
	wg.Add(1)
	// executes the async method
	go printSomething("Something", &wg)
	// and wait for it
	wg.Wait()
	// closes the writter
	_ = w.Close()
	// read the content written in the output
	result, _ := io.ReadAll(r)
	// cast the result
	output := string(result)
	// rollback the standard output
	os.Stdout = stdout
	//assertions
	if !strings.Contains(output, "Something") {
		t.Error("Expected to find 'Something', but it is not there")
	}
}
