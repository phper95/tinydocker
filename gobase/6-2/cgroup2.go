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

const (
	cgroupRoot  = "/sys/fs/cgroup"
	memoryMax   = "memory.max"   // 内存限制文件
	cpuMax      = "cpu.max"      // CPU限制文件
	cgroupProcs = "cgroup.procs" // 进程管理文件
	memoryLimit = "100m"         // 限制内存为 100MB
	cpuLimit    = "10000 100000" // 限制CPU为 10% (10000us/100000us周期)
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

// 父进程逻辑
func runParentProcess() {
	// 创建新的 cgroup 路径
	cgroupPath := filepath.Join(cgroupRoot, "memdemo")
	if err := os.Mkdir(cgroupPath, 0755); err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating cgroup: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(cgroupPath) // 清理 cgroup

	// 设置内存限制
	if err := os.WriteFile(filepath.Join(cgroupPath, memoryMax), []byte(memoryLimit), 0644); err != nil {
		fmt.Printf("Error setting memory limit: %v\n", err)
		os.Exit(1)
	}

	// 设置CPU限制
	if err := os.WriteFile(filepath.Join(cgroupPath, cpuMax), []byte(cpuLimit), 0644); err != nil {
		fmt.Printf("Error setting CPU limit: %v\n", err)
		os.Exit(1)
	}

	// 禁用swap
	if err := os.WriteFile(
		filepath.Join(cgroupPath, "memory.swap.max"),
		[]byte("0"),
		0644,
	); err != nil {
		fmt.Printf("Failed to disable swap: %v\n", err)
		os.Exit(1)
	}

	// 准备子进程命令
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID, // 创建新命名空间
	}

	// 启动子进程
	if err := cmd.Start(); err != nil {
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

	// 将子进程 PID 加入 cgroup
	pid := cmd.Process.Pid
	if err := os.WriteFile(filepath.Join(cgroupPath, cgroupProcs), []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("Error adding process to cgroup: %v\n", err)
		os.Exit(1)
	}

	// 等待子进程退出
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Child process error: %+v\n", err)
		return
	}
	fmt.Println("Parent: Child process exited")
}

// 子进程逻辑
func runChildProcess() {
	fmt.Printf("Child: PID = %d\n", os.Getpid())

	// 尝试分配200MB内存 (超过设置的100MB限制)
	var chunks [][]byte
	for i := 0; i < 20; i++ {
		// 每次分配10MB
		chunk := make([]byte, 10*1024*1024)
		for j := range chunk {
			chunk[j] = byte(j % 256) // 确保内存被实际分配
		}
		chunks = append(chunks, chunk)
		fmt.Printf("Allocated %dMB\n", (i+1)*10)
		time.Sleep(1 * time.Second) // 缩短等待时间

		// 添加延迟，让OOM killer有时间触发
		if i >= 9 { // 在分配100MB后暂停一下
			fmt.Println("Pausing to allow OOM killer to trigger...")
		}
	}

	// 如果内存分配成功，执行CPU密集型任务
	fmt.Println("Starting CPU intensive task for 10 seconds...")
	start := time.Now()
	count := 0
	for time.Since(start) < 10*time.Second {
		// 执行一些计算
		for i := 0; i < 1000000; i++ {
			count += i % 3
		}
	}
	fmt.Printf("CPU task completed. Iterations: %d\n", count)
	fmt.Println("Child: Memory allocated and CPU task completed successfully")
}
