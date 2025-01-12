package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	fmt.Println(sum(1, 2, 3, 4, 5)) // 输出：15
}
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}
