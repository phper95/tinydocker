package cgroups

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"path/filepath"
	"strconv"
)

/**
docker run -d --cpu-period=100000 --cpu-quota=250000 my_container
-cpu-period：调度周期（单位：微秒），默认 100000。
--cpu-quota：每个周期内可使用的最大 CPU 时间（单位：微秒）。
250000 / 100000 = 2.5 表示最多使用 2.5 个 CPU 核心。



docker run -d --cpus="1.5" my_container
--cpus 参数用于对容器的 CPU 使用进行硬性限制,基于 Linux Cgroups 的 cpu.cfs_period_us 和 cpu.cfs_quota_us 机制实现。
其中：
cfs_period_us 表示调度周期（单位：微秒）。
cfs_quota_us 表示该容器在每个周期内最多可以使用的 CPU 时间（单位：微秒）。
当你使用 --cpus="2.5" 时，Docker 会将其转换为：
cfs_period_us = 100000（即 100ms）
cfs_quota_us = 150000（即 1.5 * 100000）
这表示容器在一个 100ms 的周期内最多使用 150ms 的 CPU 时间（即 1.5 个 CPU 核心）。


docker run -d --name container1 --cpu-shares=512 my_image
docker run -d --name container2 --cpu-shares=1024 my_image

设置 CPU 使用的相对权重,（相对值，仅在资源竞争时生效）。
默认值：每个容器的 cpu-shares 默认值为 1024。
相对权重：数值越大，表示该容器在 CPU 竞争时能获得更多的 CPU 时间。
不设上限：如果没有其他容器竞争，一个容器可以使用全部空闲 CPU 资源，即使它的 cpu-shares 值较低。
当系统 CPU 资源紧张时，container2 将获得 container1 两倍的 CPU 时间。
如果只有一个容器运行，则它可以使用所有可用 CPU 资源，不受 cpu-shares 限制。


docker run -d --cpuset-cpus="0,1" my_container
容器只能运行在编号为 0 和 1 的 CPU 核心上。用于绑核、隔离性能敏感服务。



如果需要精确限制 CPU 使用量，使用 --cpus="X.X" 是最直观且推荐的方式。
如果只是希望调整容器之间的资源优先级，使用 --cpu-shares=X。
如果有性能隔离需求（如数据库、AI 计算等），使用 --cpuset-cpus=X。
docker中cgroup默认路径： /sys/fs/cgroup/system.slice/docker-<container-id>.scope/

*/

const (
	MemoryMax   = "memory.max"
	CpuMax      = "cpu.max"
	CgroupProcs = "cgroup.procs"
	CgroupRoot  = "/sys/fs/cgroup"
)

type CGroupManager struct {
	path string
}

func NewCGroupManager(name string) *CGroupManager {
	cgroupPath := filepath.Join(CgroupRoot, name)
	if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cgroupPath, 0755); err != nil {
			logger.Error("Error creating cgroup: %v", err)
			os.Exit(1)
		}
	}

	return &CGroupManager{path: cgroupPath}
}

func (c *CGroupManager) Apply(pid int) error {
	pidStr := strconv.Itoa(pid)
	if err := os.WriteFile(filepath.Join(c.path, CgroupProcs), []byte(pidStr), 0644); err != nil {
		logger.Error("Error applying pid to cgroup: %v", err)
		return err
	}
	return nil
}

func (c *CGroupManager) SetMemoryLimit(memoryLimit string) error {
	memPath := filepath.Join(c.path, MemoryMax)
	if err := os.WriteFile(memPath, []byte(memoryLimit), 0644); err != nil {
		logger.Error("Error setting memory limit: %v", err)
		return err
	}
	// 禁用swap
	swapPath := filepath.Join(c.path, "memory.swap.max")
	if err := os.WriteFile(swapPath, []byte("0"), 0644); err != nil {
		logger.Error("Failed to disable swap: %v", err)
		return err
	}

	return nil
}
func (c *CGroupManager) SetCPULimit(cpusStr string) error {
	// 设置CPU限制
	cpuQuota, cpuPeriod, err := ParseCPUs(cpusStr)
	if err != nil {
		logger.Error("Error parsing cpus: %v", err)
		return err
	}
	cpuMaxPath := filepath.Join(c.path, CpuMax)
	cpuLimit := fmt.Sprintf("%d %d", cpuQuota, cpuPeriod)
	if err := os.WriteFile(cpuMaxPath, []byte(cpuLimit), 0644); err != nil {
		return err
	}
	return nil
}

func (c *CGroupManager) Cleanup() error {
	err := os.RemoveAll(c.path)
	if err != nil {
		logger.Error("Error cleaning up cgroup: %v", err)
		return err
	}
	return nil
}

// ParseCPUs 将 --cpus 字符串解析为 quota 和 period
func ParseCPUs(cpusStr string) (int, int, error) {
	// 支持浮点数输入，例如 "1.5"
	cpusFloat, err := strconv.ParseFloat(cpusStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid cpus value: %s", cpusStr)
	}

	if cpusFloat <= 0 {
		return 0, 0, fmt.Errorf("cpus must be greater than 0")
	}

	// 固定周期为 100ms（CGroups 推荐值）
	const period = 100000 // microseconds
	quota := int(cpusFloat * float64(period))

	// 防止溢出或非法值
	if quota <= 0 {
		return 0, 0, fmt.Errorf("calculated quota is invalid: %d", quota)
	}

	return quota, period, nil
}
