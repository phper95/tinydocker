package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	var ap *[3]int

	//会导致panic
	for i, p := range ap {
		fmt.Println(i, p)
	}

	//舍弃for range的第二个值
	//不会导致panic
	for i, _ := range ap {
		fmt.Println(i)
	}
}
