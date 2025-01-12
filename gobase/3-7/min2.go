package main

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Config struct {
	CPU    int
	Memory int
}

type Option func(*Config)

func WithCPU(cpu int) Option {
	return func(c *Config) {
		c.CPU = cpu
	}
}

func WithMemory(memory int) Option {
	return func(c *Config) {
		c.Memory = memory
	}
}

func NewContainer(opts ...Option) *Config {
	config := &Config{}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

func main() {
	// 调用
	container := NewContainer(WithCPU(2), WithMemory(512))
	fmt.Printf("CPU: %d, Memory: %d\n", container.CPU, container.Memory) // 输出：CPU: 2, Memory: 512
}
