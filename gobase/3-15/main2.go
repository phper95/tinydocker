package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go worker(ctx, wg)
	cancel()
	wg.Wait()
}

func worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				log.Println("canceled")
			} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Println("timed out")
			}
			return
		default:
			log.Println("working")
			time.Sleep(time.Second)
		}
	}
}
