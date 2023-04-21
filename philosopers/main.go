package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []Philosopher{
	{name: "Plato", rightFork: 0, leftFork: 1},
	{name: "Socrates", rightFork: 1, leftFork: 2},
	{name: "Aristotle", rightFork: 2, leftFork: 3},
	{name: "Pascal", rightFork: 3, leftFork: 4},
	{name: "Locke", rightFork: 4, leftFork: 0},
}

var finished []string
var finishing sync.Mutex
var hunger = 2 // how many times does a person eat?
var eatTime = 500 * time.Millisecond
var thinkTime = 1 * time.Second

func main() {
	// print out a welcome message
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("---------------------------")
	fmt.Println("The table is empty.")

	//start the meal
	dine()

	// print out a finished message
	fmt.Println("The table is empty.")
	fmt.Printf("Order finished: %s.\n", strings.Join(finished, ", "))
}

func dine() {

	eated := &sync.WaitGroup{}
	eated.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)

	for index := range philosophers {
		forks[index] = &sync.Mutex{}
	}

	for _, philosopher := range philosophers {
		go eat(philosopher, eated, seated, forks)
	}

	eated.Wait()
}

func eat(philosopher Philosopher, eated *sync.WaitGroup, seated *sync.WaitGroup, forks map[int]*sync.Mutex) {
	defer eated.Done()

	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done()
	seated.Wait()

	for i := hunger; i > 0; i-- {

		if philosopher.leftFork < philosopher.rightFork {
			forks[philosopher.leftFork].Lock()
			forks[philosopher.rightFork].Lock()
		} else {
			forks[philosopher.rightFork].Lock()
			forks[philosopher.leftFork].Lock()
		}

		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()
		fmt.Printf("\t%s put down the forks.\n", philosopher.name)
	}

	fmt.Printf("%s is satisfied and left the table.\n", philosopher.name)

	finishing.Lock()
	finished = append(finished, philosopher.name)
	finishing.Unlock()
}
