package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var lock sync.Mutex
	count := 0
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock()
			count++
			lock.Unlock()
			//fmt.Println("Hello, world!")
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
