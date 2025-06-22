package container

import (
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"syscall"
)

func Run(args cli.Args, enableTTY bool) error {
	cmd := args.Get(0)
	argv := []string{"init", cmd}
	logger.Debug("args is ", argv)
	initCmd := exec.Command("/proc/self/exe", argv...)

	initCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC,
	}

	// 设置交互模式
	if enableTTY {
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		initCmd.Stdin = os.Stdin
	}

	if err := initCmd.Start(); err != nil {
		logger.Error("Failed to start container process error: ", err)
		return err
	}

	// 创建CGroup
	// cg := cgroups.NewCGroupManager(enum.AppName)
	// defer cg.Cleanup()
	// cg.SetCPULimit(50) // 限制CPU为50%
	// cg.Apply(cmd.Process.Pid)

	return initCmd.Wait()

}

// 容器内部执行的初始化函数

func mountProc() {
	// syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	target := "/proc"
	moutflags := syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	if err := syscall.Mount("proc", target, "proc", uintptr(moutflags), ""); err != nil {
		logger.Error("Failed to mount /proc: ", err)
	}
}
