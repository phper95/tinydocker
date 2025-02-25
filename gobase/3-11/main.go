package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		ch1 <- 1
	}()

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		ch2 <- 2
	}()

	go func() {
		wg.Wait()
		close(ch1)
		close(ch2)
		close(done)
	}()

	// for {
	// 	select {
	// 	case val := <-ch1:
	// 		println("ch1 received", val)
	// 	case val := <-ch2:
	// 		println("ch2 received", val)
	// 	case <-done:
	// 		return
	// 	}
	// }
	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for val := range ch1 {
			log.Println(val)
		}
	}()

	for val := range ch2 {
		log.Println(val)
	}
	wg1.Wait()
}
