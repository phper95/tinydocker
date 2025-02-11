package main

import (
	"bufio"
	"os"
)

func main() {
	//覆盖写
	//file, err := os.Create("output.txt")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//_, err = file.WriteString("Hello, docker!")
	//if err != nil {
	//	panic(err)
	//}

	//err := os.WriteFile("output.txt", []byte("Hello, docker!"), 0644)
	//if err != nil {
	//	panic(err)
	//}

	//追加写
	//file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//_, err = file.WriteString("Hello, world!\n")
	//if err != nil {
	//	panic(err)
	//}
	//通过buffer进行写入
	//创建文件
	file, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//创建writer对象
	writer := bufio.NewWriter(file)
	//写入数据
	for i := 0; i < 10; i++ {
		_, err = writer.WriteString("Hello, world!\n")
		if err != nil {
			panic(err)
		}
	}
	//刷新缓冲区
	writer.Flush()
}
