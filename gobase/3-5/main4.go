package main

import (
	"log"
	"os"
)

func main() {
	for i := 0; i < 10000000; i++ {
		// do something with file
		go processFile()

	}
}

func processFile() {
	file, err := os.Open("file.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
}
