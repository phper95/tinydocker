package main

import (
	"log"
	"sync"
)

func main() {
	m := make(map[int]int)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			m[key] = key
		}(i)
	}
	wg.Wait()
	log.Println("done")
}
