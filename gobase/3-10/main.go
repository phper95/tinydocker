package main

import (
	"log"
	"sync"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 任务完成时减 1
	log.Printf("Worker %d is starting\n", id)
	// 模拟任务
	log.Printf("Worker %d is done\n", id)
}

func main() {
	var wg sync.WaitGroup

	// 启动 3 个 Goroutine
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 增加计数器
		go worker(i, &wg)
	}

	wg.Wait() // 等待所有 Goroutine 完成
	log.Println("All workers finished")
}
