package container

import (
	"errors"
	"github.com/phper95/tinydocker/filesys"
	"github.com/phper95/tinydocker/pkg/logger"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func init() {
	logger.SetLevel(logger.DEBUG)
	logger.SetOutput(os.Stdout)
	logger.SetIncludeTrace(true)
}

func InitContainerProcess() error {
	// 从管道中读取用户传递过来的命令参数
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
	cmdArgs := strings.Split(msgStr, " ")
	if len(cmdArgs) == 0 {
		logger.Error("InitContainerProcess user cmd is empty", cmdArgs)
		return errors.New("InitContainerProcess user cmd is empty")
	}
	logger.Debug("InitContainerProcess user cmd: ", cmdArgs)
	err = filesys.Mount()
	if err != nil {
		logger.Error("Failed to mount proc: ", err)
		return err
	}

	// 在系统的PATH中寻找命令的绝对路径(因为用户可能只输入了命令名而没有输入绝对路径)
	path, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		logger.Error("exec look path error %v", err)
		return err
	}
	logger.Debug("InitContainerProcess user cmd abs path: ", path)

	// init进程读取了父进程传递过来的参数，在子进程内执行，完成了将用户指定命令传递给子进程的操作
	err = syscall.Exec(path, cmdArgs, os.Environ())
	if err != nil {
		logger.Error("Failed to exec command: ", err)
	}

	return nil
}
