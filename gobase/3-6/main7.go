package main

import "os"

func main() {
	//创建一个文件夹
	err := os.Mkdir("test", 0777)
	if err != nil {
		panic(err)
	}

}
