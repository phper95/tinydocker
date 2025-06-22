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
	// - CLONE_NEWUTS 设置新的 UTS namespace（允许设置主机名）
	// - CLONE_NEWPID 设置新的 PID namespace（容器内看到的是独立的进程ID）
	// - CLONE_NEWNS 设置新的 Mount namespace（允许挂载/卸载文件系统而不影响宿主机）
	// - CLONE_NEWIPC 设置新的 IPC namespace（隔离进程间通信）

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

func MountProc() error {
	target := "/proc"
	// MS_NODEV 禁止访问设备文件。在该文件系统中，任何字符或块设备文件都将无法被打开。防止容器内通过设备文件访问宿主机硬件资源。
	// MS_NOEXEC 禁止执行可执行文件。防止在该文件系统中运行任何程序（如 /proc 中一般不会执行程序）。
	// MS_NOSUID 禁止设置 set-user-ID 或 set-group-ID 权限。防止利用 SUID/SGID 提权，提高安全性。
	moutflags := syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	if err := syscall.Mount("proc", target, "proc", uintptr(moutflags), ""); err != nil {
		logger.Error("Failed to mount /proc: ", err)
		return err
	}
	return nil
}
