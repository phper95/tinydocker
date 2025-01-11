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
	// 模拟一些 Goroutine 的泄露
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go doTask(wg)
	}
	log.Println("Waiting for all tasks to complete...")
	wg.Wait()
}

func doTask(wg *sync.WaitGroup) {
	defer wg.Done()
	//log.Println("doing task")
	//下载文件
	time.Sleep(time.Millisecond * 200)
	//解析文件
	time.Sleep(time.Millisecond * 200)
	//入库
	time.Sleep(time.Millisecond * 200)
}
