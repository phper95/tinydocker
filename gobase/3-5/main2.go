package main

import (
	"fmt"
	"time"
)

func main() {
	x := 10
	go func() {
		x = 20
	}()
	for i := 0; i < 10; i++ {
		fmt.Println(x)
		time.Sleep(time.Second)
	}
}
