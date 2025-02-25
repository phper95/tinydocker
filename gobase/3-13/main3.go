package main

import (
	"log"
	"sync"
)

func main() {
	var count int
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, 1)
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- struct{}{}
			count++
			<-ch
		}()
	}
	wg.Wait()
	log.Println(count)
}
