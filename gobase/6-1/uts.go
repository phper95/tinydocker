package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Println("UTS Namespace Demo")

	// 显示当前主机名
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return
	}
	fmt.Printf("Current hostname: %s\n", hostname)

	// 设置子进程的命名空间标志
	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, // 创建新的UTS namespace
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

	// 在子进程中修改主机名
	// 注意：这需要在子进程的shell中手动执行，例如：
	// hostname new-hostname

	fmt.Println("Child process started with new UTS namespace")
	fmt.Println("Try running 'hostname new-hostname' in the child shell")
	fmt.Println("Then exit and check the hostname in the parent namespace")

	// 等待子进程结束
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}
}
