package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const numberOfPizzas = 10

var pizzasMade, pizzasFailed, total int
var random *rand.Rand

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func pizzeria(worker *Producer) {
	// keep track of wich pizza we are making
	var i = 0
	// run forever or until we receive a quit notification
	// try to make pizzas
	for {
		currentPizza := makePizza(i)

		i = currentPizza.pizzaNumber

		select {

		case worker.data <- *currentPizza:

		case quit := <-worker.quit:
			close(worker.data)
			close(quit)
			return
		}
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {

	pizzaNumber++

	if pizzaNumber > numberOfPizzas {
		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
		}
	}

	fmt.Printf("Received order #%d\n", pizzaNumber)

	delay := random.Intn(5) + 1

	fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
	// delay for a bit
	time.Sleep(time.Duration(delay) * time.Second)

	successful := random.Intn(12) + 1
	success := false
	message := ""

	if successful <= 4 {
		pizzasFailed++
	} else {
		pizzasMade++
	}
	total++

	if successful <= 2 {
		message = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
	} else if successful <= 4 {
		message = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
	} else {
		success = true
		message = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
	}

	result := PizzaOrder{
		pizzaNumber: pizzaNumber,
		message:     message,
		success:     success,
	}
	return &result
}

func main() {
	// seed the random number generator
	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	// print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJob)

	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= numberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas...")
			pizzaJob.Close()
		}
	}

	// print out the ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")
	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day...")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day!")
	default:
		color.Green("It was a great day!")
	}
}
