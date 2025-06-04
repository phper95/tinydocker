package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	ch := make(chan int)
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("timeout")
			return
		case i := <-ch:
			fmt.Println("received ", i)
		}
	}()
	time.Sleep(200 * time.Millisecond)
}
