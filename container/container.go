package container

import (
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Run(args cli.Args, enableTTY bool) error {
	read, write, err := os.Pipe()
	if err != nil {
		logger.Error(" Failed to create pipe:", err)
		return err
	}
	cmd := exec.Command("/proc/self/exe", "init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC,
	}

	// 设置交互模式
	if enableTTY {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
	cmd.ExtraFiles = []*os.File{read}
	if err := cmd.Start(); err != nil {
		logger.Error("Failed to start container process error: ", err)
		return err
	}

	// 创建CGroup
	// cg := cgroups.NewCGroupManager(enum.AppName)
	// defer cg.Cleanup()
	// cg.SetCPULimit(50) // 限制CPU为50%
	// cg.Apply(cmd.Process.Pid)
	command := strings.Join(args, " ")
	logger.Debug("command all is", command)
	write.WriteString(command)
	write.Close()

	return cmd.Wait()
}

// 容器内部执行的初始化函数

func mountProc() {
	target := "/proc"
	moutflags := syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	if err := syscall.Mount("proc", target, "proc", uintptr(moutflags), ""); err != nil {
		logger.Error("Failed to mount /proc: ", err)
	}
}
