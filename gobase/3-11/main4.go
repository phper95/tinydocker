package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	taskQueue := make(chan int, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			select {
			case taskQueue <- i:
				log.Println("send", i)
			default:
				log.Println("任务队列已满，任务", i, "丢弃")
			}
			time.Sleep(time.Millisecond * 200)
		}
		close(taskQueue)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case task, ok := <-taskQueue:
				if !ok {
					log.Println("任务队列已关闭,消费结束")
					return
				}
				log.Println("consume", task)
				time.Sleep(time.Millisecond * 500)
			case <-time.After(time.Second * 1):
				log.Println("暂无新任务，等待1s后继续")
			}
		}
	}()

	wg.Wait()
	log.Println("main exit")

}
