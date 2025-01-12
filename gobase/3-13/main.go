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
	wg := &sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++
		}()
	}
	wg.Wait()
	log.Printf("count: %d", count)
}
