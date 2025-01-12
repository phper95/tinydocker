package main

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 创建高优先级和低优先级任务队列
	highPriority := make(chan string, 10) // 增大缓冲区以模拟更多任务
	lowPriority := make(chan string, 10)

	// 模拟任务生成
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()
	go func() { lowPriority <- "Low Priority Task " }()
	go func() { highPriority <- "High Priority Task " }()

	// 优先处理高优先级任务的循环
	for {
		select {
		case task := <-highPriority: // 先处理高优先级任务
			log.Println("Processing High Priority:", task)
		default: // 若没有高优先级任务，再尝试处理低优先级任务
			select {
			case task := <-lowPriority:
				log.Println("Processing Low Priority:", task)
				//case <-time.After(time.Minute): // 避免低优先级阻塞
				//	log.Println("No tasks available")
			}
		}
	}
}
