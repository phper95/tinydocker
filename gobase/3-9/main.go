package main

import "fmt"

func main() {
	//单向只读channel
	//readOnlyChannel := make(<-chan int)

	//单向只写channel
	//writeOnlyChannel := make(chan<- int)

	//双向可读可写channel
	readWriteChannel := make(chan int)

	var c chan int
	fmt.Println(readWriteChannel, c) // nil
}
