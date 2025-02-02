package main

import (
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	x := 10
	go func() {
		x = 20
	}()
	for i := 0; i < 2; i++ {
		log.Println("x:", x)
		time.Sleep(time.Second)
	}
}
