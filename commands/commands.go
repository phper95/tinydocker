package commands

import (
	"errors"

	"github.com/phper95/tinydocker/container"
	"github.com/phper95/tinydocker/image"
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
		&cli.StringFlag{
			Name:  "name",
			Usage: "Assign a name to the container",
		},
		&cli.BoolFlag{
			Name:  "it",
			Usage: "Interactive mode with pseudo-TTY",
		},
		&cli.BoolFlag{
			Name:  "d",
			Usage: "Run container in detached mode (background)",
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
		name := ctx.String("name")
		if name == "" {
			return errors.New("Container name cannot be empty")
		}
		enableTTY := ctx.Bool("it")
		detach := ctx.Bool("d")

		if enableTTY && detach {
			logger.Error("-it and -d cannot be used together")
			return errors.New("-it and -d cannot be used together")
		}

		memoryLimit := ctx.String("m")
		cpuLimit := ctx.String("cpus")
		volume := ctx.String("v")
		logger.Debug("enableTTY:", enableTTY, "detach:", detach,
			"memoryLimit:", memoryLimit, "cpuLimit:", cpuLimit, "volume:", volume)
		err := container.Run(args, name, enableTTY, detach, memoryLimit, cpuLimit, volume)
		if err != nil {
			logger.Error("Run container error:", err)
		}
		return err
	},
}

// docker export imageName
var ExportCommand = cli.Command{
	Name:  "export",
	Usage: "Package the current running container into a tar file (docker export -o <tarfile> <imageName>)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "o",
			Usage: "Output file name for the tar file (default is container.tar)",
		},
	},
	Action: func(ctx *cli.Context) error {
		output := ctx.String("o")
		if output == "" {
			output = "container.tar"
		}
		if err := image.Export(output); err != nil {
			logger.Error("export error: ", err)
			return err
		}
		return nil
	},
}

// docker ps
var PsCommand = cli.Command{
	Name:  "ps",
	Usage: "List containers",
	Action: func(ctx *cli.Context) error {
		return container.PrintContainersInfo()
	},
}

// docker logs
var LogsCommand = cli.Command{
	Name:  "logs",
	Usage: "Fetch the logs of a container",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "f",
			Usage: "Follow log output",
		},
	},
	Action: func(ctx *cli.Context) error {
		containerID := ctx.Args().First()
		if containerID == "" {
			return errors.New("container name cannot be empty")
		}
		follow := ctx.Bool("f")
		return container.PrintContainerLogs(containerID, follow)
	},
}
