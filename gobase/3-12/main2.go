package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

const MaxConcurrent = 5

func init() {
	go PrintGroutineCount(time.Second)
}
func main() {
	sem := make(chan struct{}, MaxConcurrent)
	wg := sync.WaitGroup{}
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(taskID int) {
			defer wg.Done()
			log.Println("Task", taskID)
			time.Sleep(time.Millisecond * 800)
			<-sem
		}(i)
	}
	wg.Wait()
	fmt.Println("All tasks completed.")
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
