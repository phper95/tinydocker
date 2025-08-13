package main

import (
	"log"
	"os"
)

//  `os.Getenv(key string) string`: 返回指定键的环境变量值，如果不存在则返回空字符串。
//  `os.LookupEnv(key string) (string, bool)`: 返回指定键的环境变量值和一个布尔值，表示该环境变量是否存在。

func main() {
	// 读取环境变量
	value := os.Getenv("DB_HOST")
	log.Println("value:", value)
	// 检查环境变量是否存在
	value, exists := os.LookupEnv("DB_HOST")
	log.Println("value:", value, "exists:", exists)
	// 设置环境变量
	os.Setenv("DB_HOST", "localhost")

	// 获取所有环境变量
	env := os.Environ()
	log.Println("env:", env)
}
