package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// 模拟长时间文件处理，支持上下文取消
func processFile(ctx context.Context, filePath string, resultChan chan string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("❌ 文件打开失败：", err)
		return err
	}
	defer func() {
		file.Close()
		close(resultChan)
		log.Println("processFile exit")
	}()

	// 模拟分块处理文件（共处理3个块，每个块耗时1秒）
	for i := 0; i < 3; i++ {
		fmt.Printf("🔥 正在处理第 %d 个块", i+1)
		select {
		case <-ctx.Done(): // 检测取消信号
			log.Printf("取消处理，正在清理...（已处理 %d 个块）\n", i)
			return ctx.Err()
		default:
			// 模拟处理文件块耗时
			time.Sleep(1 * time.Second)
			log.Printf("✅ 成功处理第 %d 个块\n", i+1)
			resultChan <- fmt.Sprintf("第 %d 块的处理结果", i+1)
		}
	}

	return nil
}

func main() {
	doRequst()
	fmt.Println("doRequst exit")
	// 等待退出信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func doRequst() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	resultChan := make(chan string, 3)

	// 执行文件处理
	go processFile(ctx, "example.txt", resultChan)
	// 汇总处理结果
	content := ""

	// 等待处理结果或超时
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				log.Println("🎉 文件处理完成，最终处理结果：", content)
				return content, nil
			} else {
				log.Println("处理结果：", result)
				content += result
			}
		case <-ctx.Done():
			log.Println("⏰ 处理超时，已取消操作")
			return content, ctx.Err()
		}
	}

}
