package container

import (
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"syscall"
)

func init() {
	logger.SetLevel(logger.DEBUG)
	logger.SetOutput(os.Stdout)
	logger.SetIncludeTrace(true)
}

func InitContainerProcess(cmd string) error {
	mountProc()
	// 设置主机名
	err := syscall.Sethostname([]byte(enum.AppName))
	if err != nil {
		logger.Error("Failed to set hostname: ", err)
	}
	argv := []string{cmd}
	// init进程读取了父进程传递过来的参数，在子进程内执行，完成了将用户指定命令传递给子进程的操作
	err = syscall.Exec(cmd, argv, os.Environ())
	if err != nil {
		logger.Error("Failed to exec command: ", err)
	}
	return err
}
