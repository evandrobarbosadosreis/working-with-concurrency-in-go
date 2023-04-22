package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func listeToChannel(ch chan int) {
	for i := range ch {
		fmt.Println("Got", i, "from channel")
		time.Sleep(1 * time.Second) // simulate doing a lot of work
	}
	wg.Done()
}

func main() {

	channel := make(chan int, 10)

	wg.Add(1)

	go listeToChannel(channel)

	for i := 1; i <= 10; i++ {
		fmt.Println("Sending", i, "to channel...")
		channel <- i
		fmt.Println("Sent", i, "to channel!")
		fmt.Println()
	}

	close(channel)

	wg.Wait()
	fmt.Println("-----------------")
	fmt.Println("All done!")
}
