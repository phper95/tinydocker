package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		defer func() {
			ch <- 1
		}()
		// do something here
		time.Sleep(1 * time.Second)
		fmt.Println("done")

	}()
	<-ch

}
