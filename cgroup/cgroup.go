package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	cgroupPath  = "/sys/fs/cgroup"
	memoryLimit = "100000000" // 100MB
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	fmt.Println("UTS Namespace + Cgroups Memory Limit Demo", os.Args[0])

	// 显示当前主机名
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return
	}
	fmt.Printf("Current hostname: %s\n", hostname)

	// 设置子进程的命名空间标志
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// 设置子进程的输入/输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting command:", err)
	}
	log.Printf("Child process PID: %d\n", cmd.Process.Pid)

	// 启动子进程前设置cgroup
	if err := setupCgroupV2(cmd.Process.Pid); err != nil {
		fmt.Printf("Error setting up cgroup: %v\n", err)
		return
	}
	defer cleanupCgroup(cmd.Process.Pid)

	// 启动子进程
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	fmt.Println("Child process started with:")
	fmt.Println("1. New UTS namespace (isolated hostname)")
	fmt.Println("2. Memory limit of 100MB via cgroups")
	fmt.Println("Try these commands in the child shell:")
	fmt.Println("  hostname new-hostname  # Change hostname (only in this namespace)")
	fmt.Println("  stress --vm 1 --vm-bytes 150M  # Test memory limit (should be killed)")

	// 等待子进程结束
	if _, err := cmd.Process.Wait(); err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
}

func setupCgroupV2(pid int) error {
	cgroupName := fmt.Sprintf("go_demo_%d", pid)
	cgroupDir := filepath.Join(cgroupPath, cgroupName)

	// 创建 cgroup
	if err := os.MkdirAll(cgroupDir, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup dir: %v", err)
	}

	// 添加进程
	if err := ioutil.WriteFile(filepath.Join(cgroupDir, "cgroup.procs"),
		[]byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %v", err)
	}

	// 设置内存限制
	if err := ioutil.WriteFile(filepath.Join(cgroupDir, "memory.max"),
		[]byte(memoryLimit), 0644); err != nil {
		return fmt.Errorf("failed to set memory limit: %v", err)
	}

	return nil
}

func cleanupCgroup(pid int) {
	cgroupName := fmt.Sprintf("go_demo_%d", pid)
	cgroupDir := filepath.Join(cgroupPath, cgroupName)

	// 移除cgroup目录
	if err := os.RemoveAll(cgroupDir); err != nil {
		fmt.Printf("Failed to remove cgroup dir: %v\n", err)
	}
}
