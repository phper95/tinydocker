package main

import (
	"log"
	"sync"
)

func main() {
	m := make(map[int]int)
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			lock.Lock()
			m[key] = key
			lock.Unlock()
		}(i)
	}
	wg.Wait()
	log.Println("done")
}
