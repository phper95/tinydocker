package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch1 <- 1
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch2 <- 2
		}
	}()

	wg.Wait()
	close(ch1)
	close(ch2)
	for {
		select {
		case val := <-ch1:
			log.Println("ch1 received", val)
		case val := <-ch2:
			log.Println("ch2 received", val)
		}
		time.Sleep(1 * time.Second)
	}
}
