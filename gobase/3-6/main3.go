package main

import (
	"io"
	"log"
	"os"
)

func main() {
	//打开文件
	file, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//读取文件内容
	b := make([]byte, 0, 1)

	for {
		log.Println("Reading file... len(b)", len(b), "cap(b)", cap(b))
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
			log.Println("Resizing buffer... len(b)", len(b), "cap(b)", cap(b))
		}
		n, err := file.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			//文件读取完毕
			if err == io.EOF {
				err = nil
			} else {
				log.Println("Error reading file:", err)
			}
			break
		}
	}
	//输出文件内容
	log.Println("File contents:", string(b))
}
