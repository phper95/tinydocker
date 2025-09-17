package main

import "fmt"

func main() {
	defer test()()
	fmt.Println("This is the end of the program.")
}

func test() func() {
	fmt.Println("Hello, world!")
	return func() {
		fmt.Println("Goodbye, world!")
	}
}
