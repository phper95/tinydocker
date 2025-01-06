package main

import (
	"bufio"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	//file, err := os.Create("output.txt")
	//if err != nil {
	//	log.Println("Error creating file:", err)
	//	return
	//}
	//defer file.Close()
	//
	//_, err = file.WriteString("This is a new file content!")
	//if err != nil {
	//	log.Println("Error writing to file:", err)
	//}

	//追加写
	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("\nAppended content!")
	if err != nil {
		log.Println("Error appending to file:", err)
	}

	//io.WriteFile
	str := "This is a new file content!"
	err = os.WriteFile("output.txt", []byte(str), 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
	}

	// bufio.NewWriter
	file, err = os.Create("output.txt")
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	//创建writer对象
	writer := bufio.NewWriter(file)
	for i := 0; i < 10; i++ {
		_, err = writer.WriteString(string(i) + "\n")
		if err != nil {
			log.Println("Error writing to file:", err)
		}
	}
	//将缓存中的内容写入文
	writer.Flush()

}
