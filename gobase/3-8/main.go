package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	var c chan int
	var c2 chan<- int
	var c3 <-chan int
	c2 = c
	c3 = c
	//c = c2 //编译不通过
	//c = c3 //编译不通过
	//c2 = c3 //编译不通过
	fmt.Println(c2, c3)
}
