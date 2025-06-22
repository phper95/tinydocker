package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var cgroupRoot = "/sys/fs/cgroup"

const (
	memoryMax   = "memory.max"
	cpuMax      = "cpu.max"
	cgroupProcs = "cgroup.procs"
	memoryLimit = "50m"    // 50MB
	cpuQuota    = "10000"  // 10ms
	cpuPeriod   = "100000" // 100ms (10% CPU)
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Starting cgroup memory and CPU limit demo", os.Args)
	switch len(os.Args) {
	case 1:
		runParentProcess()
	case 2:
		if os.Args[1] == "child" {
			runChildProcess()
			return
		}
		fallthrough
	default:
		fmt.Println("Usage: sudo go run main.go")
		os.Exit(1)
	}
}

func runParentProcess() {
	// 创建 cgroup 路径
	cgroupPath := filepath.Join(cgroupRoot, "memdemo")

	// 确保目录存在
	if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
		if err := os.Mkdir(cgroupPath, 0755); err != nil {
			fmt.Printf("Error creating cgroup: %v\n", err)
			os.Exit(1)
		}
	}

	// 设置内存限制
	memPath := filepath.Join(cgroupPath, memoryMax)
	if err := os.WriteFile(memPath, []byte(memoryLimit), 0644); err != nil {
		fmt.Printf("Error setting memory limit: %v\n", err)
	}

	// 设置CPU限制
	cpuMaxPath := filepath.Join(cgroupPath, cpuMax)
	cpuLimit := cpuQuota + " " + cpuPeriod
	if err := os.WriteFile(cpuMaxPath, []byte(cpuLimit), 0644); err != nil {
		fmt.Printf("Error setting CPU limit: %v\n", err)
		fmt.Println("Troubleshooting:")
		fmt.Println("1. Verify CPU controller is available in the cgroup:")
		fmt.Println("   cat", filepath.Join(cgroupRoot, "cgroup.controllers"))
		fmt.Println("2. Check if cpu.max file exists:")
		fmt.Println("   ls -l", cpuMaxPath)
		fmt.Println("3. Ensure system supports cgroup v2")
		os.Exit(1)
	}

	// 禁用swap
	swapPath := filepath.Join(cgroupPath, "memory.swap.max")
	if err := os.WriteFile(swapPath, []byte("0"), 0644); err != nil {
		fmt.Printf("Failed to disable swap: %v\n", err)
	}

	// 准备子进程
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	// 启动子进程
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start child process: %v", err)
	}

	// 将子进程加入cgroup
	pid := cmd.Process.Pid
	procsPath := filepath.Join(cgroupPath, cgroupProcs)
	if err := os.WriteFile(procsPath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("Error adding process to cgroup: %v\n", err)
	}

	// 等待子进程
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() && status.Signal() == syscall.SIGKILL {
					log.Println("Child process was OOM killed (expected behavior)")
					return
				}
			}
		}
		log.Printf("Child process exited with error: %v\n", err)
	}
}

func runChildProcess() {
	fmt.Printf("Child: PID = %d\n", os.Getpid())

	// 分配少量内存 (50MB)
	for i := 0; i < 5; i++ {
		chunk := make([]byte, 10*1024*1024)
		for j := range chunk {
			chunk[j] = byte(j % 256)
		}
		fmt.Printf("Allocated %dMB\n", (i+1)*10)
		time.Sleep(500 * time.Millisecond)
	}

	// CPU 密集型任务
	fmt.Println("Starting CPU intensive task...")
	start := time.Now()
	duration := 30 * time.Second
	count := 0

	for time.Since(start) < duration {
		// 计算密集型循环
		for i := 0; i < 1000000; i++ {
			count += i % 3
		}

		// 定期报告进度
		elapsed := time.Since(start).Seconds()
		if elapsed > 0 {
			progress := elapsed / float64(duration.Seconds()) * 100
			cpuUsage := float64(count) / elapsed / 1000000
			fmt.Printf("Progress: %.1f%% | CPU Usage: %.1f Mops/s\n",
				progress,
				cpuUsage)
		}
	}

	fmt.Printf("CPU task completed. Total operations: %d\n", count)
}
