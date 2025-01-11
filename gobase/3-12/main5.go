package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	//使用channel充当锁
	var ch = make(chan struct{}, 1)
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//加锁
			ch <- struct{}{}
			count++
			<-ch
		}()
	}

	wg.Wait()
	fmt.Println("Final count:", count)
}
