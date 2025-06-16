package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	switch len(os.Args) {
	case 1:
		parentProcess()
	case 2:
		if os.Args[1] == "child" {
			childProcess()
			return
		}
		fallthrough
	default:
		fmt.Println("Usage: sudo go run main.go")
		os.Exit(1)
	}
}

func parentProcess() {
	// 创建新的UTS namespace
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start child: %v\n", err)
		os.Exit(1)
	}

	pid := cmd.Process.Pid
	fmt.Printf("Child process PID: %d\n", pid)

	// 创建cgroup目录 - 使用用户cgroup路径
	userCgroup := filepath.Join("/sys/fs/cgroup", "user.slice", "memdemo")
	if err := os.MkdirAll(userCgroup, 0755); err != nil {
		fmt.Printf("Failed to create cgroup dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(userCgroup) // 清理cgroup

	// 激活内存控制器
	if err := os.WriteFile(
		filepath.Join(userCgroup, "cgroup.subtree_control"),
		[]byte("+memory"),
		0644,
	); err != nil {
		fmt.Printf("Failed to activate memory controller: %v\n", err)
		os.Exit(1)
	}

	// 在用户cgroup下创建子cgroup
	cgroupPath := filepath.Join(userCgroup, "childcg")
	if err := os.Mkdir(cgroupPath, 0755); err != nil {
		fmt.Printf("Failed to create child cgroup: %v\n", err)
		os.Exit(1)
	}

	// 设置内存限制为100MB
	limit := "100000000" // 100MB in bytes
	if err := os.WriteFile(filepath.Join(cgroupPath, "memory.max"), []byte(limit), 0644); err != nil {
		fmt.Printf("Failed to set memory limit: %v\n", err)
		os.Exit(1)
	}

	// 添加进程到cgroup
	if err := os.WriteFile(filepath.Join(cgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("Failed to add process to cgroup: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Set memory limit to 100MB for PID %d\n", pid)

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() && status.Signal() == syscall.SIGKILL {
					fmt.Println("Child process was OOM killed (expected behavior)")
					return
				}
			}
		}
		fmt.Printf("Child process exited with error: %v\n", err)
	}
}

func childProcess() {
	fmt.Println("Child process started - allocating memory...")

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

		// 添加延迟，让OOM killer有时间触发
		if i >= 9 { // 在分配100MB后暂停一下
			fmt.Println("Pausing to allow OOM killer to trigger...")
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Println("Memory allocation completed without OOM kill")
}
