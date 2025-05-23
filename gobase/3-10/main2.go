package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {

			defer wg.Done()
			fmt.Println("Goroutine id: ", i)
			time.Sleep(time.Second)
		}(i)
	}
	wg.Wait()
	fmt.Println("Done")

}
