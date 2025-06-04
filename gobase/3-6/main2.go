package main

import (
	"io"
	"log"
	"os"
)

func main() {
	//按照指定字节来读取文件

	// 1. 打开文件
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 2. 读取指定字节
	buffer := make([]byte, 1)
	b := make([]byte, 10)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				log.Println("read file error:", err)
			}
			break
		}
		b = append(b, buffer[:n]...)
		log.Println("len(b):cap(b)", len(b), cap(b))
	}
	log.Printf("read %d str: %s", len(b), string(b))

}
