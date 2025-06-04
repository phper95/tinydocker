package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg)
	}
	wg.Wait()
	fmt.Println("Done")
}

func processFile(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Processing file: ", filename)
	time.Sleep(time.Second)
	fmt.Println("Finished processing file: ", filename)
}
