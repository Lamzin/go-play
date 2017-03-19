package main

import (
	"time"
	"fmt"
	"math/rand"
)

var (
	petrovQueue, necheporchukQueue chan item

	count int

	pValue, nValue int
)

type item struct {
	Price int
}

func main() {

	petrovQueue = make(chan item, 10)
	necheporchukQueue = make(chan item, 10)

	go ivanov()
	go petrov()
	go necheporchuk()

	for {
		time.Sleep(time.Second)
	}

}

func ivanov() {
	for {
		time.Sleep(time.Second)
		//fmt.Println("i get task")
		petrovQueue <- item{Price: rand.Intn(10) + 1}
		//fmt.Println("i complete task")
	}
}

func petrov() {
	for i := range petrovQueue {
		//fmt.Println("p get task")
		necheporchukQueue <- i
		//fmt.Println("p complete task")
	}
}

func necheporchuk() {
	for i := range necheporchukQueue {
		//fmt.Println("n get task")
		count += i.Price
		//fmt.Println("n complete task")

		fmt.Printf("Count: %d\n", count)
	}
}