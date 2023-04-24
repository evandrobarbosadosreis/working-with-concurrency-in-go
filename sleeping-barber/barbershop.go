package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	Clients         chan string
	BarberIsDone    chan bool
	Open            bool
}

func (shop *BarberShop) AddBarber(barber string) {

	shop.NumberOfBarbers++

	go func() {
		isSleeping := false

		color.Yellow("%s, goes to the waiting room to check for clients.", barber)

		for {

			// if there are no clients, the barber goes to sleep
			if len(shop.Clients) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			// aguarda receber um cliente ou fechar (Linha 68)
			client, shopOpen := <-shop.Clients

			if shopOpen {

				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}
				shop.CutHair(barber, client)

			} else {
				// shop is closed, so, send the barber home and close this goroutine
				shop.SendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) CutHair(barber, client string) {
	color.Green("%s is cuttin %s's hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
}

func (shop *BarberShop) SendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarberIsDone <- true
}

func (shop *BarberShop) CloseShopForDay() {

	shop.Open = false

	color.Cyan("Closing shop for the day.")
	close(shop.Clients) // Executa SendBarbersHome (linha 38)

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarberIsDone // Aguarda todos os barbeiros irem pra casa (Linha 63)
	}

	close(shop.BarberIsDone) // fecha o canal
	color.Green("--------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day and everyone has gone home.")
}

func (shop *BarberShop) AddClient(name string) {
	color.Green("*** %s arrives!", name)

	if shop.Open {
		select {
		case shop.Clients <- name:
			color.Yellow("%s takes a seat in the waiting room.", name)
		default:
			color.Red("The waiting room is full, so %s leaves.", name)
		}

	} else {
		color.Red("The shop is already closed, so %s leaves!", name)
	}
}
