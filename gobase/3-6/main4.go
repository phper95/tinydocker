package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	//逐行读取文件内容
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//创建reader
	reader := bufio.NewReader(file)
	var content string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("文件读取完毕")
			} else {
				log.Println(err)
			}
			break
		}
		//打印每行内容
		log.Println(line)
		content += line
	}
	log.Println("文件内容:", content)
}
