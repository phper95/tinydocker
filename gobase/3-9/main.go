package main

func main() {
	//单向只读channel
	//readOnlyChannel := make(<-chan int)
	//readOnlyChannel <- 10 // Error: cannot send to receive-only channel
	//单向只写channel
	//writeOnlyChannel := make(chan<- int)
	//<-writeOnlyChannel // Error: cannot receive from send-only channel

	//双向可读可写channel
	//readWriteChannel := make(chan int)

	//var c chan int
	//fmt.Println(readWriteChannel == nil, c == nil)

	var c chan int
	var c2 chan<- int
	var c3 <-chan int
	//c2 = c
	//c3 = c
	c = c2
	//c = c3
	c2 = c3

	// 带缓存区的 channel
	bufferedChannel := make(chan int, 10)

	// 只读且不带缓存区 channel
	unbufferedReadOnlyChannel := make(<-chan int, 4)

	// 只写且带缓存区 channel
	bufferedWriteOnlyChannel := make(chan<- int, 10)
}
