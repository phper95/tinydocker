package filesys

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"syscall"
)

func CreateOverlayFS() error {
	// 解压 busybox-rootfs.tar 到 /var/local/busybox
	tarPath := "busybox-rootfs.tar"
	extractPath := "/var/local/busybox"

	if err := extractTarGz(tarPath, extractPath); err != nil {
		logger.Error("failed to extract %s to %s: %v", tarPath, extractPath, err)
		return err
	}

	// 设置 OverlayFS 参数
	lowerDir := "/var/local/busybox"
	upperDir := "/tmp/upper"
	workDir := "/tmp/work"
	mountPoint := "/mnt/overlay"

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

	// 执行 mount 命令
	if err := syscall.Mount("overlay", mountPoint, "overlay", 0, options); err != nil {
		return fmt.Errorf("failed to mount overlayfs: %w", err)
	}

	fmt.Println("OverlayFS mounted successfully")
	return nil
}
