package container

import (
	"github.com/phper95/tinydocker/enum"
	"log"

	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func InitContainerProcess() error {
	// 0-stdin
	// 1-stdout
	// 2-stderr
	// 3-pipe
	pipe := os.NewFile(uintptr(3), "pipe")

	// block read
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Println("init read pipe error %v", err)
		return nil
	}

	msgStr := string(msg)
	cmdArray := strings.Split(msgStr, " ")
	log.Println("cmdArray is", cmdArray)
	// mountProc()

	path, err := exec.LookPath(cmdArray[0])
	log.Println("path is", path)
	if err != nil {
		log.Println("exec look path error %v", err)
		return err
	}
	return nil
	// init进程读取了父进程传递过来的参数，在子进程内执行，完成了将用户指定命令传递给子进程的操作
	syscall.Exec(path, cmdArray[0:], os.Environ())
	syscall.Sethostname([]byte(enum.AppName))
	return nil
}
