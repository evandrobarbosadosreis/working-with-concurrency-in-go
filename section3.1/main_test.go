package main

import (
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {

	var m sync.Mutex

	wg.Add(2)
	go updateMessage("Hello, world!", &m)
	go updateMessage("Hello, world!", &m)
	wg.Wait()

	if msg != "Hello, world!" {
		t.Error("You have a race condition")
	}
}
