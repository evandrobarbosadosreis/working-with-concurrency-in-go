package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	fmt.Println(s)
	defer wg.Done()
}

func main() {

	var wg sync.WaitGroup

	words := []string{
		"One",
		"Two",
		"Three",
		"Four",
		"Five",
		"Six",
		"Seven",
		"Eight",
		"Nine",
		"Ten",
	}

	wg.Add(len(words))

	for i, x := range words {
		go printSomething(fmt.Sprintf("%d : %s", i, x), &wg)
	}
	wg.Wait()
}
