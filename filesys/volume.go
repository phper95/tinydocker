package filesys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/phper95/tinydocker/pkg/logger"
)

// MountVolume 解析 -v 传入的 "hostDir:containerDir" 并完成 bind-mount
func MountVolume(volume, mountRoot string) error {
	if volume == "" {
		logger.Warn("no volume specified")
		return nil // 未指定数据卷
	}

	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid volume format %q, expected hostDir:containerDir", volume)
	}
	hostDir, containerDir := parts[0], parts[1]

	// 校验宿主机目录是否存在
	if _, err := os.Stat(hostDir); err != nil {
		return fmt.Errorf("host dir %s not exist: %w", hostDir, err)
	}

	// 在容器根目录下创建目标挂载点
	dest := filepath.Join(mountRoot, containerDir)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("create container dir %s failed: %w", dest, err)
	}

	// 绑定挂载
	// 为什么要使用绑定挂载？
	/**
	  普通挂载：通过文件系统驱动程序读取存储设备（如硬盘、网络存储）上的数据结构，并将其映射为用户可见的目录树。
	  绑定挂载：直接在 VFS（虚拟文件系统）层创建一个新的挂载点，指向已存在的文件或目录，不涉及底层文件系统的解析。
	  //绑定挂载是Linux内核提供的特殊机制，它可以：
	  //1. 可以实现跨文件系统的数据共享，是容器技术实现自身文件系统的基础。
	  //2.零拷贝开销
	  // 不复制数据，只是建立新的路径映射
	  // 访问速度与直接访问原文件相同
	*/
	if err := syscall.Mount(hostDir, dest, "none", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("bind mount %s to %s failed: %w", hostDir, dest, err)
	}
	logger.Debug("volume mounted %s --> %s", hostDir, dest)
	return nil
}

func UnmountVolume(volume, mountRoot string) error {
	if volume == "" {
		return nil
	}
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		// If format is invalid we cannot know mount point; log and ignore.
		return fmt.Errorf("invalid volume spec %q", volume)
	}
	_, containerDir := parts[0], parts[1]
	dest := filepath.Join(mountRoot, containerDir)

	if err := syscall.Unmount(dest, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount volume %s failed: %w", dest, err)
	}
	logger.Debug("Volume unmounted %s", dest)
	return nil
}
