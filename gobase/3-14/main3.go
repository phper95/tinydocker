package main

import (
	"log"
	"sync"
)

func main() {
	m := make(map[int]int)
	lock := sync.RWMutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			lock.Lock()
			m[key] = key
			lock.Unlock()
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			val, ok := m[key]
			lock.RLock()
			log.Println(val, ok)
			lock.RUnlock()
		}(i)

	}
	wg.Wait()
	log.Println("done")
}
