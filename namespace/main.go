package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// 阶段1：创建子进程命令（含Namespace声明）
	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS, // 组合隔离
	}

	// 阶段2：配置子进程环境（UTS隔离演示）
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"PS1=[新Namespace] # "} // 提示符标记

	// 阶段3：启动并验证隔离性
	if err := cmd.Run(); err != nil {
		fmt.Printf("启动失败: %v\n", err)
		os.Exit(1)
	}
}
