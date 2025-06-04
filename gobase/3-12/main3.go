package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

func init() {
	go PrintGroutineCount(time.Second)
}
func main() {
	wg := sync.WaitGroup{}
	taskQueue := make(chan func(), 5)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskQueue {
				task()
			}
		}()
	}
	for i := 0; i < 10000000; i++ {
		taskQueue <- func() {
			doTask()
		}
	}

	log.Println("Waiting for all tasks to complete...")
	wg.Wait()
}
func PrintGroutineCount(dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("Current Groutine count: %d", runtime.NumGoroutine())
		}
	}
}
func doTask() {
	// do some task here
	// 下载文件
	time.Sleep(time.Millisecond * 200)
	// 解压文件
	time.Sleep(time.Millisecond * 300)
	// 入库
	time.Sleep(time.Millisecond * 100)
}
