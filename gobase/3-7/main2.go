package main

import "fmt"

type Config struct {
	Cpu    int
	Memory int
	Disk   int
}

type Options func(*Config)

func WithCpu(cpu int) Options {
	return func(config *Config) {
		config.Cpu = cpu
	}
}

func WithMemory(memory int) Options {
	return func(config *Config) {
		config.Memory = memory
	}
}

func WithDisk(disk int) Options {
	return func(config *Config) {
		config.Disk = disk
	}
}

func NewContanier(opts ...Options) {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	fmt.Printf("cpu:%d,memory:%d\n", config.Cpu, config.Memory)
}

func main() {
	NewContanier(WithMemory(2))

}
