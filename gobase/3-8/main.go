package main

import (
	"flag"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 定义一个命令参数
	name := flag.String("name", "world", "a name to greet")

	// 解析命令行参数
	flag.Parse()

	// 输出结果
	log.Printf("Hello, %s!\n", *name)
}
