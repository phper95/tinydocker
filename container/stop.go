package container

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"path/filepath"
	"syscall"
)

// Stop stops a running container
func Stop(containerName string) error {
	// 查找容器信息
	info, err := findContainerInfo(containerName)
	if err != nil {
		return fmt.Errorf("failed to find container %s: %v", containerName, err)
	}

	// 检查容器是否正在运行
	if info.State != ContainerStateRunning {
		return fmt.Errorf("container %s is not running", containerName)
	}

	// 向容器进程发送终止信号
	pid := info.Pid
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		logger.Error("Failed to send signal to container %s: %v", containerName, err)
		return fmt.Errorf("failed to send signal to container %s: %v", containerName, err)
	}

	// 更新容器状态为已停止
	err = UpdateContainerState(info.Id, ContainerStateStopped)
	if err != nil {
		return fmt.Errorf("failed to update container %s state: %v", containerName, err)
	}

	logger.Info("Container %s stopped", containerName)
	return nil
}

// findContainerInfo finds container info by name or ID
func findContainerInfo(nameOrID string) (*Info, error) {
	// 首先尝试按ID精确匹配
	filePath := filepath.Join(DefaultContainerInfoPath, nameOrID, DefaultContainerInfoFileName)
	if _, err := os.Stat(filePath); err == nil {
		// 找到精确匹配的ID
		return ReadContainerInfo(filePath)
	}

	// 如果没有精确匹配，遍历所有容器查找匹配名称的容器
	return GetContainerInfoByName(nameOrID)
}
