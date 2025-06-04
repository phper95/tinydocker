package main

import (
	"log"
	"time"
)

func main() {
	ch := make(chan int)
	go doWork(ch)
	select {
	case msg := <-ch:
		log.Println("Received", msg)
	case <-time.After(2 * time.Second):
		log.Println("timed out")
	}
}

func doWork(ch chan int) {
	time.Sleep(3 * time.Second)
	ch <- 1
}
