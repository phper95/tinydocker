package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		//x := i
		//go func(parm int) {
		//	defer wg.Done()
		//	fmt.Println(parm)
		//}(i)
		go print(i, &wg)
	}
	wg.Wait()
	fmt.Println("Done")
}

func print(i int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if r := recover(); r != nil {
			fmt.Println("Recovered in print", r)
		}
	}()
	fmt.Println(i)
	panic("Oops")
}
