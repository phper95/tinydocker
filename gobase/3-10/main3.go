package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		//defer wg.Done()
		fmt.Println("first goroutine")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("second goroutine")
	}()

	wg.Wait()
	fmt.Println("Done")

}
