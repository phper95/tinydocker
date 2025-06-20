package commands

import (
	"context"
	"github.com/phper95/tinydocker/container"
	"github.com/urfave/cli/v3"
	"log"
)

var RunCommand = &cli.Command{
	// 命令名称
	Name:  "run",
	Usage: "Run a command in a new container",
	// 命令参数
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "Interactive mode with pseudo-TTY",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		// 获取命令参数列表
		args := cmd.Args().Slice()
		// 命令行参数校验
		if len(args) == 0 {
			log.Fatal("No command specified")
		}
		return container.Run(args, cmd.Bool("it"))
	},
}
