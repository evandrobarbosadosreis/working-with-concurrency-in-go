package main

import (
	"fmt"
	"time"
)

func server1(ch chan<- string) {
	for {
		time.Sleep(6 * time.Second)
		ch <- "This is coming from server 1"
	}
}

func server2(ch chan<- string) {
	for {
		time.Sleep(3 * time.Second)
		ch <- "This is coming from server 2"
	}
}

func main() {

	channel1 := make(chan string)
	channel2 := make(chan string)

	go server1(channel1)
	go server2(channel2)

	for {
		select {
		case s1 := <-channel1:
			fmt.Println("Case 1 says: ", s1)
		case s2 := <-channel1:
			fmt.Println("Case 2 says: ", s2)
		case s3 := <-channel2:
			fmt.Println("Case 3 says: ", s3)
		case s4 := <-channel2:
			fmt.Println("Case 4 says: ", s4)
		}
	}
}
