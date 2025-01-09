package main

import (
	"context"
	"errors"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 创建一个 2秒超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保退出时释放资源
	// 模拟一个耗时任务
	done := make(chan struct{}, 1)
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		time.Sleep(3 * time.Second)
		log.Println("任务完成！")
	}()
	select {
	case <-done:
		log.Println("任务完成！")
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("任务超时！")
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			log.Println("任务取消！")
		}

	}
}
