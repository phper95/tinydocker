package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Worker id: ", id)
	time.Sleep(time.Second * 2)
	log.Println("Worker id: ", id, "done")
}
