package filesys

import (
	"github.com/phper95/tinydocker/pkg/logger"
	"syscall"
)

func MountProc() error {
	target := "/proc"
	// MS_NODEV 禁止访问设备文件。在该文件系统中，任何字符或块设备文件都将无法被打开。防止容器内通过设备文件访问宿主机硬件资源。
	// MS_NOEXEC 禁止执行可执行文件。防止在该文件系统中运行任何程序（如 /proc 中一般不会执行程序）。
	// MS_NOSUID 禁止设置 set-user-ID 或 set-group-ID 权限。防止利用 SUID/SGID 提权，提高安全性。
	if err := syscall.Mount("proc", target, "proc", syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_NOSUID, ""); err != nil {
		logger.Error("Failed to mount /proc: ", err)
		return err
	}
	return nil
}
