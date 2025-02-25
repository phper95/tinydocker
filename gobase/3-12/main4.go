package main

import (
	"log"
	"runtime"
	"time"
)

type Task func() error

type workerPool struct {
	tasks      chan Task
	taskSize   int
	workerSize int
}

func init() {
	go PrintGroutineCount(time.Second)
}

// 初始化协程池并自动启动worker
func newWorkerPool(taskSize, workerSize int) *workerPool {
	pool := &workerPool{
		tasks:      make(chan Task, taskSize),
		taskSize:   taskSize,
		workerSize: workerSize,
	}

	// 启动worker
	pool.run()
	return pool
}
func (p *workerPool) run() {
	for i := 0; i < p.workerSize; i++ {
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("worker", workerID, "panic:", r)
				}
			}()
			for task := range p.tasks {
				if task == nil {
					continue
				}
				log.Println("worker", workerID, "start task")
				if err := task(); err != nil {
					log.Println("worker", workerID, "task error:", err)
				}
			}
		}(i)
	}
}

func (p *workerPool) Submit(task Task) {
	p.tasks <- task
}
func main() {
	pool := newWorkerPool(10, 5)

	for i := 0; i < 1000000000; i++ {
		pool.Submit(func() error {
			doTask()
			return nil
		})
	}
}
func PrintGroutineCount(dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("Current Groutine count: %d", runtime.NumGoroutine())
		}
	}
}
func doTask() {
	log.Println("do task")
	panic("do task panic")
	// do some task here
	// 下载文件
	time.Sleep(time.Millisecond * 200)
	// 解压文件
	time.Sleep(time.Millisecond * 300)
	// 入库
	time.Sleep(time.Millisecond * 100)
}
