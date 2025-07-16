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
		// cmd := ctx.Args().Get(0)
		// logger.Debug("init command:", cmd)
		return container.InitContainerProcess()
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
			Name:  "cpus",
			Usage: "CPU limit for the container (e.g., 1.5)",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "Bind mount a volume (host_dir:container_dir)",
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
		memoryLimit := ctx.String("m")
		cpuLimit := ctx.String("cpus")
		volume := ctx.String("v")
		logger.Debug("enableTTY:", enableTTY, "memoryLimit:", memoryLimit, "cpuLimit:", cpuLimit, "volume:", volume)
		err := container.Run(args, enableTTY, memoryLimit, cpuLimit, volume)
		if err != nil {
			logger.Error("Run container error:", err)
		}
		return err
	},
}
