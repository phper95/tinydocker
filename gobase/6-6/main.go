package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// 创建内存管道
	reader, writer := io.Pipe()
	var wg sync.WaitGroup

	// 启动生产者协程
	wg.Add(1)
	go func() {
		defer writer.Close() // 生产结束后关闭写入端
		defer wg.Done()

		// 模拟生产数据
		messages := []string{
			"Hello, pipe!",
			"This is a sample message.",
			"Using io.Pipe for communication.",
			"Last message, closing pipe soon...",
		}

		for _, msg := range messages {
			// 写入数据到管道
			if _, err := writer.Write([]byte(msg + "\n")); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to pipe: %v\n", err)
				return
			}
			fmt.Println("Produced:", msg)
			time.Sleep(500 * time.Millisecond) // 控制生产速度
		}
	}()

	// 启动消费者协程
	wg.Add(1)
	go func() {
		defer wg.Done()

		// 从管道读取数据
		scanner := io.TeeReader(reader, os.Stdout) // 将数据同时输出到标准输出
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, scanner); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading from pipe: %v\n", err)
			return
		}

		// 处理完整数据
		fmt.Println("\n--- Received Messages ---")
		fmt.Println(buf.String())
	}()

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("Communication completed.")
}
