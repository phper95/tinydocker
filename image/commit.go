package image

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/phper95/tinydocker/container"
	"github.com/phper95/tinydocker/pkg/logger"
)

// Commit 将正在运行的容器文件系统打成 tar 包
// imageName 为生成镜像的文件名前缀，例如 test -> test.tar
func Commit(imageName string) error {
	if imageName == "" {
		return fmt.Errorf("image name cannot be empty")
	}

	// 组合输出路径：/var/local/images/<imageName>.tar
	outputDir := "/var/local/images"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create image dir: %w", err)
	}
	dstTar := filepath.Join(outputDir, fmt.Sprintf("%s.tar", imageName))

	rootfs := container.MountPoint // /mnt/overlay
	if _, err := os.Stat(rootfs); err != nil {
		return fmt.Errorf("container rootfs not found, is container running? %w", err)
	}

	// 使用宿主 tar 命令归档
	cmd := exec.Command("tar", "-cvf", dstTar, "-C", rootfs, ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("tar error: %w, output: %s", err, string(output))
	}

	logger.Info("container committed to %s", dstTar)
	return nil
}
