package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	ch := make(chan int, 5)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go producer(ch, 1, &wg)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go consumer(ch, i+1, &wg)
	}

	wg.Wait()
}

func producer(ch chan int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i < 6; i++ {
		log.Println("producer", id, "send", i)
		ch <- i
		time.Sleep(time.Millisecond * 500)
	}
	close(ch)
}

func consumer(ch chan int, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range ch {
		log.Println("consumer", id, "receive", val)
		time.Sleep(time.Millisecond * 800)
	}
}
