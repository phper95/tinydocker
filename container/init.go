package container

import (
	"github.com/phper95/tinydocker/pkg/logger"
	"io"
	"os"
	"strings"
	"syscall"
)

func init() {
	logger.SetLevel(logger.DEBUG)
	logger.SetOutput(os.Stdout)
	logger.SetIncludeTrace(true)
}

func InitContainerProcess(cmd string) error {
	// cmdArray := readUserCommand()
	// if cmdArray == nil || len(cmdArray) == 0 {
	// 	return fmt.Errorf("run container get user command error, cmdArray is nil")
	// }
	err := Mount()
	if err != nil {
		return err
	}
	argv := []string{cmd}
	logger.Debug("InitContainerProcess user cmd: ", cmd, " argv: ", argv)

	// 在系统的PATH中寻找命令的绝对路径
	// path, err := exec.LookPath(cmdArray[0])
	// if err != nil {
	// 	logger.Error("exec look path error %v", err)
	// 	return err
	// }

	// logger.Debug("find path %s", path)

	// init进程读取了父进程传递过来的参数，在子进程内执行，完成了将用户指定命令传递给子进程的操作
	err = syscall.Exec(cmd, argv, os.Environ())
	if err != nil {
		logger.Error("Failed to exec command: ", err)
	}
	// if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
	// 	logger.Error("exec error %v", err)
	// 	return err
	// }
	return nil
}

func readUserCommand() []string {
	// 0-stdin
	// 1-stdout
	// 2-stderr
	// 3-pipe
	pipe := os.NewFile(uintptr(3), "pipe")

	// block read
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logger.Error("init read pipe error %v", err)
		return nil
	}

	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
