package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go PrintGoroutineCount(time.Second)
}

// PrintGoroutineCount 每秒打印当前 Goroutine 数量
func PrintGoroutineCount(dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
		}
	}
}

func main() {
	wg := &sync.WaitGroup{}

	taskQueue := make(chan func(), 100)

	// 开启10个协程，同时执行任务
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

		//将执行的任务放入队列中
		taskQueue <- func() {
			doTask()
		}
	}
	log.Println("Waiting for all tasks to complete...")
	wg.Wait()
}

func doTask() {
	log.Println("doing task")
	//下载文件
	time.Sleep(time.Millisecond * 200)
	//解析文件
	time.Sleep(time.Millisecond * 200)
	//入库
	time.Sleep(time.Millisecond * 200)
}
