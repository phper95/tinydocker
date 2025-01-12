package main

import (
	"context"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("任务被取消：", ctx.Err())
			return
		default:
			log.Println("正在工作...")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx)

	time.Sleep(3 * time.Second)
	cancel() // 手动取消任务
	time.Sleep(1 * time.Second)
}
