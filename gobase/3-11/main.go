package main

import (
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func producer(ch chan int, id int) {
	for i := 1; i <= 5; i++ {
		log.Printf("Producer %d produced: %d\n", id, i)
		ch <- i
		time.Sleep(time.Millisecond * 500) // 模拟生产时间
	}
	close(ch)
}

func consumer(ch chan int, id int) {
	for val := range ch {
		log.Printf("Consumer %d consumed: %d\n", id, val)
		time.Sleep(time.Millisecond * 800) // 模拟消费时间
	}
}

func main() {
	ch := make(chan int, 5) // 缓冲区大小为 5

	go producer(ch, 1)
	go consumer(ch, 1)

	time.Sleep(time.Second * 5)
}
