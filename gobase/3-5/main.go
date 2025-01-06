package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	data, err := os.ReadFile("example.txt")
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}
	log.Println("os.ReadFile:", string(data))

	//ioutilData
	ioutilData, err := ioutil.ReadFile("example.txt")
	if err != nil {
		log.Println("Error reading file using ioutil:", err)
		return
	}
	log.Println("ioutil.ReadFile:", string(ioutilData))

	file, err := os.Open("example.txt")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	//data, err = io.ReadAll(file)
	//log.Println("io.ReadAll:", string(data))

	//data, err = ioutil.ReadAll(file)
	//log.Println("ioutil.ReadAll:", string(data))

	//固定字节读取
	//buffer := make([]byte, 1)
	//b := make([]byte, 0, 1024)
	//for {
	//	n, err := file.Read(buffer)
	//	if err != nil {
	//		if err == io.EOF {
	//			err = nil
	//		} else {
	//			log.Println("Error reading file:", err)
	//		}
	//		break
	//	}
	//	b = append(b, buffer[:n]...)
	//	//buffer = make([]byte, 1)
	//}
	//log.Println("file.Read fixed byte:", string(b))

	//按照切片扩容读取
	//b := make([]byte, 0, 1)
	//for {
	//	if len(b) == cap(b) {
	//		log.Println("append before len(b)", len(b), "cap(b)", cap(b))
	//		//触发分片扩容
	//		b = append(b, 0)[:len(b)]
	//		log.Println("append after len(b)", len(b), "cap(b)", cap(b))
	//	}
	//	//log.Println("b[len(b):cap(b)]", b[len(b):cap(b)])
	//	n, err := file.Read(b[len(b):cap(b)])
	//	b = b[:len(b)+n]
	//	if err != nil {
	//		if err == io.EOF {
	//			err = nil
	//		} else {
	//			log.Println("Error reading file:", err)
	//		}
	//		break
	//	}
	//}
	//log.Println("file.Read slice expand:", string(b))

	//通过bufio缓冲区读取文件(逐行读取大文件)
	reader := bufio.NewReader(file) //建立缓冲区,将文件内容放入到缓冲区
	for {
		line, err := reader.ReadString('\n') //读取一行内容
		if err != nil {
			if err == io.EOF {
				log.Println("bufio.ReadString:", line)
			} else {
				log.Println("Error reading file:", err)
			}
			break
		}
		log.Println("bufio.ReadString:", line)
	}

	scanner := bufio.NewScanner(file)
	// 自定义 buffer size 为 2MB
	buffer := make([]byte, 2*1024*1024)
	scanner.Buffer(buffer, len(buffer))

	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error scanning file:", err)
	}

}
