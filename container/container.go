package container

import (
	"github.com/phper95/tinydocker/cgroups"
	"github.com/phper95/tinydocker/enum"
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
	cg := cgroups.NewCGroupManager(enum.AppName)
	defer cg.Cleanup()
	cg.SetMemoryLimit("100m")
	cg.SetCPULimit(10000, 100000) // 限制CPU为50%
	cg.Apply(initCmd.Process.Pid)

	return initCmd.Wait()

}

func Mount() {

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

func MountRoofs(root string) error {
	// MS_BIND 表示绑定挂载，将宿主机的目录或文件系统挂载到容器内部。
	// Linux 的某些系统调用（如 pivot_root 或 chroot）要求目标路径必须是一个挂载点
	// 挂载自身 + 启用 MS_REC 标志后，会复制所有子挂载点；
	// 这样容器后续调用 pivot_root 切换根目录时，就能看到完整的、独立的文件系统结构（包括 /proc, /sys 等子挂载）；
	// 并且这些挂载操作不会影响宿主机或其他容器。

	// 默认情况下，挂载点可能是 shared 类型，意味着在一个命名空间中挂载或卸载会影响其他命名空间。
	// 通过绑定挂载自身：
	// 实际上创建了一个“私有副本”；
	// 防止不同命名空间之间的挂载操作互相干扰；
	// 保证容器内部挂载/卸载操作只在容器内生效

	// MS_REC启用递归挂载，将当前挂载点下的所有子挂载点也进行绑定。
	// 让容器内部看到的挂载结构与宿主机一致。
	moutflags := syscall.MS_BIND | syscall.MS_REC
	if err := syscall.Mount(root, root, "bind", moutflags, ""); err != nil {
		logger.Error("Failed to mount rootfs: ", err)
		return err
	}
	return nil
}

// 该函数的作用是将一个 tmpfs 文件系统挂载到 /dev 目录。
// tmpfs 是一种基于内存的临时文件系统，常用于需要快速读写且不需要持久化的场景
func MountTmpfs() error {
	moutflags := syscall.MS_NOSUID | syscall.MS_STRICTATIME
	// 源设备 (source): 通常为 tmpfs，表示使用虚拟文件系统而非物理设备。
	// 挂载点 (target): /dev，表示将 tmpfs 挂载到容器的 /dev 目录。
	// 文件系统类型 (fstype): tmpfs，表示使用临时文件系统。

	// 挂载标志 (flags): MS_NOSUID | MS_STRICTATIME，禁止设置 set-user-ID 或 set-group-ID 权限，并严格遵守访问时间更新策略。
	// Linux 提供了其他较宽松的 atime 更新策略，例如：
	// MS_RELATIME：仅当访问时间早于 mtime 或 ctime 时才更新 atime（默认行为）。
	// MS_NOATIME：完全禁止更新 atime。
	// 这些优化减少了磁盘 I/O，但可能影响某些应用逻辑。
	// 使用 MS_STRICTATIME 明确禁用这些优化，确保每次都更新 atime。
	// 因为 tmpfs 是内存中的文件系统，其访问速度远高于磁盘。 更新 atime 带来的额外开销非常小，几乎不会影响性能。

	// 数据字段 (data): "mode=755"，指定挂载的目录模式为 755（即 rwxr-xr-x）。
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", moutflags, "mode=755"); err != nil {
		logger.Error("Failed to mount tmpfs: ", err)
		return err
	}
	return nil
}
