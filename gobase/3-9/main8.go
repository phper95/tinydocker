package main

func main() {
	//ch := make(chan int, 3)
	//close(ch)
	//close(ch)
	//log.Println(ch == nil)

	var ch chan int
	//ch = make(chan int)
	//ch <- 1
	<-ch
}
