package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Println("PID + Mount Namespace Demo")

	// 显示当前进程ID
	fmt.Printf("Parent PID: %d\n", os.Getpid())

	// 设置子进程的命名空间标志
	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, // 创建新的PID和mount namespace
	}

	// 设置子进程的输入/输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动子进程
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	fmt.Println("Child process started with new PID and mount namespace")
	fmt.Println("In the child shell:")
	fmt.Println("1. Run 'echo $$' to see PID 1")
	fmt.Println("2. Run 'mount -t proc proc /proc' to setup proper /proc")
	fmt.Println("3. Then 'ps aux' will only show namespace processes")

	// 等待子进程结束
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
}
