package main

import (
	"fmt"
	"log"
	"sync"
)

type SafeMap struct {
	data   map[string]int
	opChan chan func(map[string]int)
}

func NewSafeMap() *SafeMap {
	safeMap := &SafeMap{
		data:   make(map[string]int),
		opChan: make(chan func(map[string]int)),
	}

	go safeMap.run()
	return safeMap
}

func (sm *SafeMap) run() {
	for op := range sm.opChan {
		op(sm.data)
	}
}

func (sm *SafeMap) Set(key string, value int) {
	sm.opChan <- func(data map[string]int) {
		data[key] = value
	}
}

func (sm *SafeMap) Get(key string) (int, bool) {
	var val int
	var ok bool
	wg := sync.WaitGroup{}
	wg.Add(1)
	sm.opChan <- func(data map[string]int) {
		defer wg.Done()
		val, ok = data[key]
	}
	wg.Wait()
	return val, ok
}

func main() {
	sm := NewSafeMap()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			sm.Set(fmt.Sprintf("%d", key), key)
		}(i)
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		v, ok := sm.Get(fmt.Sprintf("%d", i))
		log.Println(v, ok)
	}

}
