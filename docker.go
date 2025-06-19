package main

import (
	"context"
	"github.com/phper95/tinydocker/commands"
	"github.com/phper95/tinydocker/enum"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	// 创建顶级命令（替代 v2 的 cli.App）
	cmd := &cli.Command{
		Name:  enum.AppName,
		Usage: "A simple container runtime",
		Commands: []*cli.Command{ // 注意是指针切片
			commands.RunCommand,
		},
	}

	// 使用 cli.Run 执行命令
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
