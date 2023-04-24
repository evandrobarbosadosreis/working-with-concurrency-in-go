package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatinCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {

	// print welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels if we need any
	clients := make(chan string, seatinCapacity)
	done := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatinCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		Clients:         clients,
		BarberIsDone:    done,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	// add barbers
	shop.AddBarber("Frank") // irá startar uma nova rotina pra cada barbeiro
	shop.AddBarber("John")
	shop.AddBarber("Susan")
	shop.AddBarber("Pam")
	shop.AddBarber("Eddy")
	shop.AddBarber("Sarah")

	// start the barbershop as gourotines
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() { // essa rotina só aguarda mantém a loja aberta e fecha no horário
		<-time.After(timeOpen) // bloqueia por 10 segundos
		shopClosing <- true    // envia notificação de "fechando" (para de enviar clientes na linha 68/69)
		shop.CloseShopForDay() // fecha
		closed <- true         // envia notificação de "fechado" (bloqueando a thread principal na linha 78)
	}()

	// add clients

	i := 1

	var random = rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() { // outra rotina que add um novo cliente à loja a cada x milissegundos
		for {
			interval := random.Int() % (2 * arrivalRate)

			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(interval)):
				shop.AddClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	// block until barbarshop is closed
	<-closed
}
