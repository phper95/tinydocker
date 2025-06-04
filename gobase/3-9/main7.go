package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	for {
		select {
		case num, ok := <-ch:
			if !ok {
				fmt.Println("channel closed", num)
			} else {
				fmt.Println("Received", num)
			}
			time.Sleep(time.Second)
		}
	}
}
