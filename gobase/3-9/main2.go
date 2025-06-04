package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		//ch <- 1
		fmt.Println("开始读取ch")
		<-ch
		fmt.Println("读取ch结束")
	}()

	go func() {
		ch <- 1
	}()

	//log.Println(<-ch)

	time.Sleep(1 * time.Second)
}
