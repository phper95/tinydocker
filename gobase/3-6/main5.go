package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	//打开文件
	f, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//读取文件内容
	scanner := bufio.NewScanner(f)
	buffer := make([]byte, 64*1024)
	scanner.Buffer(buffer, len(buffer))
	var content string
	for scanner.Scan() {
		log.Println(scanner.Text())
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	log.Println("Content:", content)
}
