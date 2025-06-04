package main

import (
	"log"
	"sync"
	"sync/atomic"
)

func main() {
	var count int64
	wg := sync.WaitGroup{}
	// lock := sync.Mutex{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// lock.Lock()
			// defer lock.Unlock()
			// count++
			atomic.AddInt64(&count, 1)
			// lock.Unlock()
		}()
	}
	wg.Wait()
	log.Println(count)

}
