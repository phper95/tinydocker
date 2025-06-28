package main

import (
	"fmt"
	"os"
)

func main() {
	// 读取环境变量
	value := os.Getenv("PATH")

	if value == "" {
		fmt.Println("环境变量 PATH 未设置")
	} else {
		fmt.Printf("环境变量 PATH 的值为: %s\n", value)
	}

	value, exists := os.LookupEnv("PATH")

	if !exists {
		fmt.Println("环境变量 PATH 未设置")
	} else {
		fmt.Printf("环境变量 PATH 的值为: %s\n", value)
	}
}
