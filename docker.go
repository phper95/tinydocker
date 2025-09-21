package main

import (
	"github.com/phper95/tinydocker/pkg/db"
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
	log.Println("tinydocker start")
	if !isInitProcess() {
		InitBoltDB()
		defer func() {
			err := db.GetBoltDBClient(db.DefaultBoltDBClientName).Close()
			if err != nil {
				logger.Error("close bolt db error", err)
			}
		}()
	}

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
		commands.NetworkCommand,
	}

	// 使用 cli.Run 执行命令
	if err := app.Run(os.Args); err != nil {
		logger.Error("app run error", err)
		panic(err)
	}
}

// 判断是否为 init 进程
func isInitProcess() bool {
	return len(os.Args) > 1 && os.Args[1] == "init"
}

func InitBoltDB() {
	err := db.InitBoltDBClient(db.DefaultBoltDBClientName, enum.DefaultNetworkDBPath)
	if err != nil {
		logger.Error("init bolt db error", err)
		panic(err)
	}
	err = db.GetBoltDBClient(db.DefaultBoltDBClientName).CreateBucketIfNotExists(enum.DefaultNetworkTable)
	if err != nil {
		logger.Error("create network table error", err)
	}
	err = db.GetBoltDBClient(db.DefaultBoltDBClientName).CreateBucketIfNotExists(enum.AllocatedIPKey)
	if err != nil {
		logger.Error("create allocated ip table error", err)
	}
	log.Println("init bolt db finished", db.DefaultBoltDBClientName)
}
