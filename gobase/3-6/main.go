package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	//data, err := os.ReadFile("example.txt")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println(string(data))

	//ioutil.ReadFile()
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	//ioutil.ReadAll()
}
