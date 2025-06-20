package container

import (
	"github.com/phper95/tinydocker/container/process"
	"github.com/phper95/tinydocker/enum"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/phper95/tinydocker/container/cgroups"
)

func Run(command []string, it bool) error {
	cmd := exec.Command("/proc/self/exe", "child")
	cmd.Args = append(cmd.Args, command...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC,
	}

	// 设置交互模式
	if it {
		process.AttachTerminal(cmd)
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}

	if err := cmd.Start(); err != nil {
		log.Println("Failed to start container process: %v", err)
		return err
	}

	// 创建CGroup
	cg := cgroups.NewCGroupManager(enum.AppName)
	defer cg.Cleanup()
	cg.SetCPULimit(50) // 限制CPU为50%
	cg.Apply(cmd.Process.Pid)

	return cmd.Wait()
}

// 容器内部执行的初始化函数
func InitContainer() {
	mountProc()
	syscall.Sethostname([]byte(enum.AppName))
}

func mountProc() {
	target := "/proc"
	moutflags := syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	if err := syscall.Mount("proc", target, "proc", uintptr(moutflags), ""); err != nil {
		log.Fatalf("Failed to mount /proc: %v", err)
	}
}
