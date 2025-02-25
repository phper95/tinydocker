package main

import (
	"log"
	"sync"
	"time"
)

func main() {

	wg := sync.WaitGroup{}
	tasks := make(chan int, 10)

	// 启动3个worker
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, tasks, &wg)
	}

	// 分发任务
	for i := 1; i <= 10; i++ {
		tasks <- i
	}
	close(tasks)
	wg.Wait()
	log.Println("All tasks are done.")
}

func worker(id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		log.Println("Worker ", id, " processing task ", task)
		time.Sleep(time.Second)
	}
}
