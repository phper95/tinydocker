package main

import (
	"context"
	"fmt"
	"time"
)

func task(ctx context.Context, name string, duration time.Duration) {
	select {
	case <-time.After(duration):
		fmt.Println(name, "完成")
	case <-ctx.Done():
		fmt.Println(name, "取消：", ctx.Err())
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go task(ctx, "任务1", 1*time.Second)
	go task(ctx, "任务2", 3*time.Second)

	time.Sleep(4 * time.Second)
}
