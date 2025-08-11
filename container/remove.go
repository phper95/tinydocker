package container

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"path/filepath"
)

// Remove removes a container
func Remove(containerName string, force bool) error {
	// 查找容器信息
	targetInfo, err := GetContainerInfoByName(containerName)
	if err != nil {
		logger.Error("failed to get container info: %v", err)
		return err
	}

	if targetInfo == nil {
		return fmt.Errorf("container %s not found", containerName)
	}

	// 如果容器正在运行且没有使用force参数，则返回错误
	if targetInfo.State == ContainerStateRunning && !force {
		return fmt.Errorf("cannot remove running container %s, use -f to force remove", containerName)
	}

	// 如果容器正在运行且使用了force参数，则先停止容器
	if targetInfo.State == ContainerStateRunning && force {
		err := Stop(containerName)
		if err != nil {
			return fmt.Errorf("failed to stop container %s: %v", containerName, err)
		}
	}

	// 删除容器目录
	containerDir := filepath.Join(DefaultContainerInfoPath, targetInfo.Id)
	err = os.RemoveAll(containerDir)
	if err != nil {
		return fmt.Errorf("failed to remove container directory %s: %v", containerDir, err)
	}

	logger.Info("Container %s removed", containerName)
	return nil
}
