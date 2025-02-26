package main

import (
	"log"
	"sync"
)

func main() {
	// m := make(map[int]int)
	// m[1] = 2
	// m[3] = 4
	// m[0] = 0
	//
	// v, ok := m[0]
	// log.Println(v, ok)

	var safeMap sync.Map
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			safeMap.Store(key, key)
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			val, ok := safeMap.Load(key)
			log.Println(val, ok)
		}(i)
		wg.Wait()
	}
}
