package container

import (
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
	argv := []string{cmd}
	// init进程读取了父进程传递过来的参数，在子进程内执行，完成了将用户指定命令传递给子进程的操作
	syscall.Exec(cmd, argv, os.Environ())
	// syscall.Sethostname([]byte(enum.AppName))
	return nil
}
