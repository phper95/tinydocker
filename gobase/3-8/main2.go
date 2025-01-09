package main

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {

	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	log.Println(<-ch)
}
