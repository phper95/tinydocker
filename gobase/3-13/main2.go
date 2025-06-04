package main

import (
	"log"
	"sync"
)

func main() {
	ch := make(chan int, 1000)
	var count int
	wg := sync.WaitGroup{}
	// lock := sync.Mutex{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- 1
		}()

	}
	wgConsumer := sync.WaitGroup{}
	wgConsumer.Add(1)
	go func() {
		defer wgConsumer.Done()
		// 相当于将并发操作收拢到一个channel中使用一个goroutine进行处理，从而避免了并发操作导致的错误
		for val := range ch {
			count += val
		}
	}()

	wg.Wait()
	close(ch)
	wgConsumer.Wait()
	log.Println(count)

}
