package main

import (
	"context"
	"fmt"
	"time"
)

// fetchResult 模拟一个带有超时的远程调用
func fetchResult(ctx context.Context) (string, error) {
	resultChan := make(chan string)

	go func() {
		time.Sleep(3 * time.Second) // 模拟耗时操作
		resultChan <- "success"
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultChan:
		return res, nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err := fetchResult(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Result:", result)
}
