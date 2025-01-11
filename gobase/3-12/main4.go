package main

import (
	"log"
	"sync"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	//展示计数在并发环境中的错误用法
	count := 0
	wgProducer := &sync.WaitGroup{}
	wgConsumer := &sync.WaitGroup{}
	//使用chnnel实现
	ch := make(chan int, 1000)

	wgConsumer.Add(1)
	//相当于将并发操作收拢到一个channel中使用一个goroutine进行处理，从而避免了并发操作导致的错误
	go func() {
		defer wgConsumer.Done()
		for i := range ch {
			count += i
		}
	}()

	for i := 0; i < 10000; i++ {
		wgProducer.Add(1)
		go func() {
			defer wgProducer.Done()
			ch <- 1
		}()
	}
	wgProducer.Wait()
	close(ch)
	wgConsumer.Wait()
	log.Printf("count: %d", count)
}
