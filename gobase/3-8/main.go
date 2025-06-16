package main

import (
	"flag"
	"fmt"
)

func main() {
	// 定义一个命令参数
	name := flag.String("name", "default", "name")
	flag.Parse()
	fmt.Println(*name)
}
