package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// asyncTask 模拟一个复杂的耗时任务
func doTask(ctx context.Context, resultChan chan<- string) {
	defer close(resultChan) // 确保通道在任务结束时关闭

	taskSteps := 5 // 模拟任务需要执行的步骤数
	for i := 1; i <= taskSteps; i++ {
		select {
		case <-ctx.Done(): // 检查取消信号
			fmt.Printf("Task stopped at step %d: %v\n", i, ctx.Err())
			return
		default:
			// 模拟每一步的耗时操作
			processTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
			fmt.Printf("Processing step %d, will take %v...\n", i, processTime)
			time.Sleep(processTime) // 模拟操作时间
		}
	}

	// 如果任务完成
	resultChan <- "Task successfully completed"
}

func main() {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子

	// 创建一个 context，设置超时时间为 3 秒
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultChan := make(chan string)

	// 启动异步任务
	go doTask(ctx, resultChan)

	// 等待任务完成或超时
	select {
	case result := <-resultChan:
		if result == "" {
			log.Println("Task stopped unexpectedly")
			return
		}
		fmt.Println("Result:", result)
	case <-ctx.Done():
		fmt.Println("Operation timed out:", ctx.Err())

	}
	time.Sleep(10 * time.Second)
}
