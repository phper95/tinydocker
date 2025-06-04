package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// 使用 unshare 命令创建一个新的 PID namespace 并运行 sh
	cmd := exec.Command("unshare", "--pid", "sh", "-c", "echo 'In new PID namespace'; sleep 10 & exec sh")

	// 设置标准输入输出为当前进程的标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}
}
