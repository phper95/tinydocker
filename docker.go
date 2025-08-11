package main

import (
	"log"
	"os"

	"github.com/phper95/tinydocker/commands"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
)

func init() {
	logger.SetLevel(logger.DEBUG)
	logger.SetOutput(os.Stdout)
	logger.SetIncludeTrace(true)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	logger.Debug("tinydocker start")
	app := cli.NewApp()
	app.Name = enum.AppName
	app.Usage = "A simple container runtime"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		commands.InitCommand,
		commands.RunCommand,
		commands.ExportCommand,
		commands.PsCommand,
		commands.LogsCommand,
		commands.ExecCommand,
		commands.ExecContainerCommand,
		commands.StopCommand,
		commands.RemoveCommand,
	}

	// 使用 cli.Run 执行命令
	if err := app.Run(os.Args); err != nil {
		logger.Error("app run error", err)
	}
}