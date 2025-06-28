package container

import (
	"fmt"
	"github.com/phper95/tinydocker/cgroups"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Run(args cli.Args, enableTTY bool, memoryLimit, cpuLimit string) error {
	logger.Debug("Run  args: ", args)

	// initCmdArgs := []string{"init"}
	// 将Run命令的参数传递给init命令
	// if len(initCmdArgs) > 0 {
	// 	initCmdArgs = append(initCmdArgs, args...)
	// }

	initCmd, write, err := NewInitProcess(enableTTY, memoryLimit, cpuLimit)
	if err != nil {
		logger.Error("Failed to create init process error: ", err)
		return err
	}
	logger.Debug("Container process started with pid: ", initCmd.Process.Pid)

	// 将管道写入端传递给init命令
	err = SendInitCommand(args, write)
	if err != nil {
		logger.Error("Failed to send init command error: ", err)
		return err
	}
	// 等待容器退出
	if enableTTY {
		return initCmd.Wait()
	}
	return nil

}

func NewInitProcess(enableTTY bool, memoryLimit, cpuLimit string) (*exec.Cmd, *os.File, error) {

	read, write, err := os.Pipe()
	if err != nil {
		logger.Error("Failed to create pipe error: ", err)
		return nil, nil, err
	}

	initCmd := exec.Command("/proc/self/exe", "init")
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

	// 传入管道文件读取端句柄，外带此句柄去创建子进程
	initCmd.ExtraFiles = []*os.File{read}

	// 设置交互模式
	if enableTTY {
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		initCmd.Stdin = os.Stdin
	}

	if err := initCmd.Start(); err != nil {
		logger.Error("Failed to start container process error: ", err)
		return initCmd, write, err
	}

	// 创建CGroup
	cg := cgroups.NewCGroupManager(enum.AppName)
	// 设置内存限制
	if memoryLimit != "" {
		err := cg.SetMemoryLimit(memoryLimit)
		if err != nil {
			logger.Error("Failed to set memory limit error: ", err)
			return initCmd, write, err
		}

	}

	// 设置CPU限制
	if cpuLimit != "" {
		err := cg.SetCPULimit(cpuLimit) // 限制CPU为50%
		if err != nil {
			logger.Error("Failed to set cpu limit error: ", err)
			return initCmd, write, err
		}

	}

	// 应用CGroup
	err = cg.Apply(initCmd.Process.Pid)
	if err != nil {
		logger.Error("Failed to apply cgroup error: ", err)
		return initCmd, write, err
	}

	logger.Debug("NewInitProcess Container process started with pid: ", initCmd.Process.Pid)
	return initCmd, write, nil
}

func SendInitCommand(cmd []string, write *os.File) error {
	logger.Debug("send init command: %v", cmd)
	defer write.Close()
	command := strings.Join(cmd, " ")
	logger.Debug("command all is [ %s ]", command)
	if _, err := write.WriteString(command); err != nil {
		return fmt.Errorf("send init command [%s] error:%v", command, err)
	}
	return nil
}
