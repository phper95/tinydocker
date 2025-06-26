package container

import (
	"fmt"
	"github.com/phper95/tinydocker/cgroups"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func Run(args cli.Args, enableTTY bool, memoryLimit, cpuLimit string) error {
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
	// 设置内存限制
	if memoryLimit != "" {
		err := cg.SetMemoryLimit(memoryLimit)
		if err != nil {
			logger.Error("Failed to set memory limit error: ", err)
			return err
		}

	}

	// 设置CPU限制
	if cpuLimit != "" {
		err := cg.SetCPULimit(cpuLimit) // 限制CPU为50%
		if err != nil {
			logger.Error("Failed to set cpu limit error: ", err)
			return err
		}

	}

	// 应用CGroup
	err := cg.Apply(initCmd.Process.Pid)
	if err != nil {
		logger.Error("Failed to apply cgroup error: ", err)
		return err
	}

	// 等待容器退出
	return initCmd.Wait()

}

// GenerateCPULimit 根据指定的 CPU 使用百分比和周期生成 "quota period" 字符串
// 例如：percent=10, period=100000 微秒（100ms） => 返回 "10000 100000"
func GenerateCPULimit(percent int, period uint64) (string, error) {
	if percent < 0 || percent > 100 {
		return "", fmt.Errorf("percent must be between 0 and 100")
	}
	if period == 0 {
		period = 100000 // 默认使用 100ms 周期
	}

	quota := (int64(percent) * int64(period)) / 100
	return fmt.Sprintf("%d %d", quota, period), nil
}

func Mount() {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Error("Get current location error %v", err)
		return
	}
	logger.Debug("Current location is %s", pwd)
	MountProc()
	MountPivotRoot(pwd)
	MountTmpfs()
}

// 容器内部执行的初始化函数

func MountProc() error {
	// syscall.Mount(source string, target string, fstype string, flags uintptr, data string)
	// 1. source string
	// 含义：要挂载的源设备或目录。
	// 示例：
	// "proc"：表示虚拟文件系统（如 /proc）；
	// "/dev/sda1"：表示物理设备；
	// "tmpfs"：表示内存文件系统类型；
	// root：当前目录或指定路径作为绑定挂载的源

	// 2. target string
	// 含义：挂载点路径，即将 source 挂载到哪个目录。
	// 注意： // 必须是一个存在的目录；
	// 如果是绑定挂载（bind mount），则目标路径也必须是一个有效的路径。
	// 这里是挂载到 /proc。

	// 3. fstype string
	// 含义：文件系统类型。
	// 常见值：
	// "proc"：虚拟进程信息文件系统；
	// "tmpfs"：基于内存的临时文件系统；
	// "bind"：绑定挂载，将一个已存在的目录挂载到另一个位置；
	// "ext4"、"xfs" 等：实际磁盘文件系统。
	// 示例：MountTmpfs 使用 "tmpfs"。

	// 4. flags uintptr
	// 含义：挂载标志位，控制挂载行为。
	// 常用标志（可组合使用 |）：
	// MS_BIND：绑定挂载，复制一个已有的挂载点到新位置；
	// MS_RDONLY：只读挂载；
	// MS_NODEV：不允许访问设备文件；
	// MS_NOEXEC：禁止执行可执行文件；
	// MS_NOSUID：禁止 SUID 和 SGID 权限生效；
	// MS_REC：递归操作，适用于绑定挂载时复制所有子挂载；
	// MS_STRICTATIME / MS_RELATIME / MS_NOATIME：控制访问时间更新策略；
	// MS_PRIVATE / MS_SHARED / MS_SLAVE：控制命名空间中挂载传播行为。

	// 5. data string
	// 含义：传递给文件系统的额外选项或参数。
	// 格式：通常是逗号分隔的键值对字符串。
	// 示例：
	// "mode=755"：设置 tmpfs 挂载目录权限为 rwxr-xr-x；
	// "size=100m"：限制 tmpfs 大小；
	// ""：某些情况不需要参数时传空字符串。

	// /proc 是一个虚拟文件系统，提供进程和内核信息
	// 容器需要查看自己的进程信息（如 PID、内存使用等），而不是宿主机上的全局信息。
	// 通过挂载私有的 /proc，实现进程视图的隔离
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
	if err := syscall.Mount(root, root, "bind", uintptr(moutflags), ""); err != nil {
		logger.Error("Failed to mount rootfs: ", err)
		return err
	}
	return nil
}

// 将当前工作目录作为容器的新根目录（pivot_root）
// 实现文件系统的隔离。容器只能看到其根目录下的文件结构，无法访问宿主机的其他文件。
// 防止容器逃逸到宿主机文件系统
// 先绑定挂载自身目录（MountRoofs），确保子挂载点也被复制。
// 创建 .pivot_root 临时目录用于过渡。
// 使用 syscall.PivotRoot 将当前目录设为新的根目录。
// 最后卸载旧的根并清理临时目录。
func MountPivotRoot(root string) error {
	// 挂载根文件系统
	err := MountRoofs(root)
	if err != nil {
		logger.Error("Failed to mount rootfs: ", err)
		return err
	}

	// 创建一个临时目录 .pivot_root，用于在切换根目录时作为旧根的挂载点
	// .开头对用户隐藏，防止用户误操作
	pivotDir := filepath.Join(root, ".pivot_root")
	if err = os.Mkdir(pivotDir, 0755); err != nil {
		logger.Error("Failed to create pivot dir: ", err)
		return err
	}

	// 将当前进程的根文件系统切换到新的根目录 root
	// 旧的根目录会被挂载到 pivotDir 上
	if syscall.PivotRoot(root, pivotDir); err != nil {
		logger.Error("Failed to pivot root: ", err)
		return err
	}
	// 修改当前工作目录到新根目录
	if err = syscall.Chdir("/"); err != nil {
		logger.Error("Failed to chdir: ", err)
		return err
	}
	pivotDir = filepath.Join("/", ".pivot_root")
	// 解除绑定挂载
	// syscall.MNT_DETACH Linux 系统调用中用于卸载挂载点的一个标志，
	// 其作用是将指定的挂载点从文件系统中分离（detach），而不影响其他引用该挂载点的位置。
	if err = syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		logger.Error("Failed to unmount pivot dir: ", err)
		return err
	}
	if err = os.Remove(pivotDir); err != nil {
		logger.Error("Failed to remove pivot dir: ", err)
		return err
	}
	return nil
}

// 该函数的作用是将一个 tmpfs 文件系统挂载到 /dev 目录。
// tmpfs 是一种基于内存的临时文件系统，常用于需要快速读写且不需要持久化的场景
// 容器需要基本的设备节点（如 /dev/null, /dev/zero 等）来运行程序。
// 使用 tmpfs 可以动态生成这些设备节点，并且是临时的，重启后不会保留。
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
	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", uintptr(moutflags), "mode=755"); err != nil {
		logger.Error("Failed to mount tmpfs: ", err)
		return err
	}
	return nil
}
