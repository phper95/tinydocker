package main

import (
	"fmt"
	"log"
	"sync"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	m := make(map[int]int)
	var wg sync.WaitGroup

	// 启动多个 goroutine 并发写 Map
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			m[key] = key // 并发写入
		}(i)
	}

	wg.Wait()

	fmt.Println("done")
}
