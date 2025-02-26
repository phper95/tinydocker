package main

import (
	"context"
	"errors"
	"log"
	"runtime"
	"time"
)

func init() {
	go PrintGroutineCount(time.Second)
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
func main() {
	f1()
	log.Println("main done")
	select {}
}

func f1() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		defer func() {
			close(done)
			// done <- struct{}{}
		}()
		time.Sleep(time.Second * 3)
		log.Println("done")
	}()

	select {
	case <-done:
		log.Println("finished")
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("timed out")
		} else if errors.Is(ctx.Err(), context.Canceled) {
			log.Println("canceled")
		}
		return
	}
}
