package main

import (
	"fmt"
	"strings"
)

func shout(ping <-chan string, pong chan<- string) {
	for {
		s := <-ping

		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}
}

func main() {

	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and press ENTER (enter Q to quit)")

	for {

		fmt.Print("-> ")
		var userInput string
		fmt.Scan(&userInput)

		if strings.ToLower(userInput) == "q" {
			break
		}

		ping <- userInput
		response := <-pong
		fmt.Println("->", response)
	}

	fmt.Println("All done. Closing channels and going home.")
	close(ping)
	close(pong)
}
