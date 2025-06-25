package commands

import (
	"errors"
	"github.com/phper95/tinydocker/container"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
)

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		logger.Debug("init command args:", ctx.Args())
		cmd := ctx.Args().Get(0)
		logger.Debug("init command:", cmd)
		return container.InitContainerProcess(cmd)
	},
}
var RunCommand = cli.Command{
	// 命令名称
	Name:  "run",
	Usage: "Run a command in a new container",
	// 命令参数
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "Interactive mode with pseudo-TTY",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "Memory limit for the container (e.g., 512m, 1g)",
		},
		&cli.StringFlag{
			Name:  "c",
			Usage: "CPU limit for the container (e.g., 1, 2)",
		},
		&cli.StringFlag{
			Name:  "p",
			Usage: "Port mapping (e.g., 8080:80)",
		},
	},
	Action: func(ctx *cli.Context) error {
		// 获取命令参数列表
		args := ctx.Args()
		logger.Debug("args:", args)
		// 命令行参数校验
		if len(args) == 0 {
			logger.Error("No command specified")
			return errors.New("No command specified")
		}
		enableTTY := ctx.Bool("it")
		err := container.Run(args, enableTTY)
		if err != nil {
			logger.Error("Run container error:", err)
		}
		return err
	},
}
