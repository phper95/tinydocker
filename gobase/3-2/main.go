package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	slice := []int{0, 1, 2, 3}
	m := make(map[int]*int)
	for key, val := range slice {
		log.Printf("val: %d; addr: %p", val, &val)
		m[key] = &val
	}
	for k, v := range m {
		fmt.Println(k, "->", *v)
	}
}
