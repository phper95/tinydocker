package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan struct{})
	go func() {
		defer wg.Done()
		log.Println("starting goroutine")
		time.Sleep(time.Second * 2)
		log.Println("goroutine done")
	}()
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		log.Println("done")
	case <-time.After(time.Second * 5):
		log.Println("timed out")
	}

}
