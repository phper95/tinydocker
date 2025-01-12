package main

import (
	"log"
	"runtime"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go PrintGoroutineCount(time.Second)
}

type Task func() error

type WorkerPool struct {
	tasks      chan Task
	taskSize   int
	workerSize int
}

// 初始化协程池并自动启动
func NewWorkerPool(workerSize int, taskSize int) *WorkerPool {
	pool := &WorkerPool{
		tasks:      make(chan Task, taskSize), // 使用缓冲通道来限制任务数量
		workerSize: workerSize,
	}
	pool.run()
	return pool
}

// 启动协程池
func (p *WorkerPool) run() {
	for i := 0; i < p.workerSize; i++ {
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Worker %d panic: %v\n", workerID, r)
				}
			}()
			for task := range p.tasks {
				if task == nil {
					continue
				}
				log.Printf("Worker %d processing task\n", workerID)
				err := task()
				if err != nil {
					log.Printf("Worker %d task error: %v\n", workerID, err)
				}
			}
		}(i)
	}
}

// 提交任务
func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}

func main() {
	pool := NewWorkerPool(2, 100) // 初始化协程池，最多 3 个并发任务

	// 提交任务
	for i := 0; i < 1000000; i++ {
		//index := i
		pool.Submit(func() error {
			//log.Printf("Task %d started\n", index)
			doTask()
			//log.Printf("Task %d finished\n", index)
			return nil
		})
	}
}
func doTask() {
	log.Println("doing task")
	//panic("task panic")
	//下载文件
	time.Sleep(time.Millisecond * 200)
	//解析文件
	time.Sleep(time.Millisecond * 200)
	//入库
	time.Sleep(time.Millisecond * 200)
}
func PrintGoroutineCount(dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
		}
	}
}
