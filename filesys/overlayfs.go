package filesys

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/file"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"path"
	"syscall"
)

func CreateOverlayFS(busyboxDir, mountPoint, tarPath string) error {
	// 设置 OverlayFS 相关文件夹
	lowerDir := busyboxDir
	upperDir := path.Join(path.Dir(busyboxDir), "upper")
	workDir := path.Join(path.Dir(busyboxDir), "work")
	// 创建 upper 和 work 目录
	err := os.MkdirAll(upperDir, 0755)
	if err != nil {
		logger.Error("failed to create %s: %v", upperDir, err)
		return err
	}
	err = os.MkdirAll(workDir, 0755)
	if err != nil {
		logger.Error("failed to create %s: %v", workDir, err)
		return err
	}

	// 解压 busybox-rootfs.tar 到 /var/local/busybox
	if !file.IsDir(busyboxDir) {
		if err := file.ExtractTarGz(tarPath, busyboxDir); err != nil {
			logger.Error("failed to extract %s to %s: %v", tarPath, busyboxDir, err)
			return err
		}
	}

	// 挂载 OverlayFS
	if err := MountOverlayFS(lowerDir, upperDir, workDir, mountPoint); err != nil {
		logger.Error("failed to mount overlayfs: %v", err)
		return err
	}
	return nil
}

func MountOverlayFS(lowerDir, upperDir, workDir, mountPoint string) error {
	// 创建工作目录和挂载点
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return err
	}

	// 构建 overlayfs 挂载选项
	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)

	// 执行 mount 命令（挂载标志=0表示默认的挂载选项，不启用任何特殊的挂载标志）
	// MS_RDONLY  // 只读挂载
	// MS_NOSUID  // 忽略 set-user-ID 和 set-group-ID 位
	// MS_NODEV  // 不允许访问设备文件
	// MS_NOEXEC  // 不允许执行程序
	// MS_REMOUNT  // 重新挂载一个已存在的挂载点
	if err := syscall.Mount("overlay", mountPoint, "overlay", 0, options); err != nil {
		return fmt.Errorf("failed to mount overlayfs: %w", err)
	}

	logger.Debug("OverlayFS mounted successfully")
	return nil
}

func UnmountOverlayFS(busyboxDir, mountPoint string) error {
	// 卸载 OverlayFS
	if err := syscall.Unmount(mountPoint, 0); err != nil {
		logger.Error("failed to unmount overlayfs: %v", err)
		return fmt.Errorf("failed to unmount overlayfs: %w", err)
	}
	upperDir := path.Join(path.Dir(busyboxDir), "upper")
	workDir := path.Join(path.Dir(busyboxDir), "work")
	err := os.RemoveAll(upperDir)
	if err != nil {
		logger.Error("failed to remove %s: %v", upperDir, err)
	}
	err = os.RemoveAll(workDir)
	if err != nil {
		logger.Error("failed to remove %s: %v", workDir, err)
	}

	logger.Debug("OverlayFS unmounted successfully")
	return err
}
