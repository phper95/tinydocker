package main

import (
	"log"
	"os"
)

func main() {
	//创建目录
	//err := os.Mkdir("test", 0777)
	//if err != nil {
	//	panic(err)
	//}
	//err := os.MkdirAll("test/test1/test2", 0777)
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("创建目录成功")

	//删除文件
	//err := os.Remove("test/test.txt")
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("删除文件成功")
	//删除目录
	//err := os.RemoveAll("test/")
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("删除目录成功")

	//移动文件
	//err := os.Rename("test/test1/test3.txt", "test/test3.txt")
	//if err != nil {
	//	panic(err)
	//}
	//log.Println("移动文件成功")

	//移动目录
	err := os.Rename("test/test1/test2", "test/test2")
	if err != nil {
		panic(err)
	}
	log.Println("移动目录成功")
}
