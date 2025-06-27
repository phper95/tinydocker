package cgroups

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"path/filepath"
	"strconv"
)

/**
内存限制
docker run --memory="512m" --memory-swap="1g" my_image

CPU 资源限制
docker run -d --cpu-period=100000 --cpu-quota=250000 my_container
-cpu-period：调度周期（单位：微秒），默认 100000。
--cpu-quota：每个周期内可使用的最大 CPU 时间（单位：微秒）。
250000 / 100000 = 2.5 表示最多使用 2.5 个 CPU 核心。



docker run -d --cpus="1.5" my_container
--cpus 参数用于对容器的 CPU 使用进行硬性限制,基于 Linux Cgroups 的 cpu.cfs_period_us 和 cpu.cfs_quota_us 机制实现。
其中：
cfs_period_us 表示调度周期（单位：微秒）。
cfs_quota_us 表示该容器在每个周期内最多可以使用的 CPU 时间（单位：微秒）。
当你使用 --cpus="1.5" 时，Docker 会将其转换为：
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
	MemoryMax   = "memory.max"   // 内存限制配置文件，用于设置cgroup的内存上限
	CpuMax      = "cpu.max"      // CPU限制配置文件，用于设置cgroup的CPU使用上限
	CgroupProcs = "cgroup.procs" // cgroup进程列表文件，用于将进程加入指定的cgroup
	CgroupRoot  = "/sys/fs/cgroup" // cgroup挂载根目录，是Linux系统中管理控制组的默认路径
)

type CGroupManager struct {
	path string
}

// NewCGroupManager 创建一个新的 CGroupManager 对象
//
// 参数:
//     name - cgroup 的名称
//
// 返回值:
//     *CGroupManager - 指向新创建的 CGroupManager 对象的指针
//
// 注意:
//     1. 如果指定的 cgroup 路径不存在，则会尝试创建该路径
//     2. 如果创建路径失败，则会记录错误日志并退出程序
func NewCGroupManager(name string) *CGroupManager {
	// 拼接cgroup路径
	cgroupPath := filepath.Join(CgroupRoot, name)

	// 检查cgroup路径是否存在
	if _, err := os.Stat(cgroupPath); os.IsNotExist(err) {
		// 如果路径不存在，尝试创建路径
		if err := os.MkdirAll(cgroupPath, 0755); err != nil {
			// 记录错误日志
			logger.Error("Error creating cgroup: %v", err)
			// 退出程序
			os.Exit(1)
		}
	}

	// 返回CGroupManager对象
	return &CGroupManager{path: cgroupPath}
}


// Apply 将给定的进程ID（pid）加入到 cgroup 中。
//
// 参数：
//     pid：要加入 cgroup 的进程ID。
//
// 返回值：
//     如果成功，则返回 nil；否则返回错误信息。
func (c *CGroupManager) Apply(pid int) error {
	// 将 pid 转换为字符串
	pidStr := strconv.Itoa(pid)

	// 将 pid 写入到 cgroup 的 procs 文件中
	if err := os.WriteFile(filepath.Join(c.path, CgroupProcs), []byte(pidStr), 0644); err != nil {
		// 如果写入失败，记录错误日志
		logger.Error("Error applying pid to cgroup: %v", err)
		return err
	}

	// 写入成功，返回 nil
	return nil
}


// SetMemoryLimit 为CGroup设置内存限制
//
// 参数:
//     memoryLimit: 设置的内存限制值
//
// 返回值:
//     如果设置成功，返回nil；如果设置失败，返回错误信息
func (c *CGroupManager) SetMemoryLimit(memoryLimit string) error {
	// 拼接路径
	memPath := filepath.Join(c.path, MemoryMax)
	// 写入内存限制
	if err := os.WriteFile(memPath, []byte(memoryLimit), 0644); err != nil {
		// 记录错误日志
		logger.Error("Error setting memory limit: %v", err)
		// 返回错误
		return err
	}
	// 禁用swap
	// 拼接swap路径
	swapPath := filepath.Join(c.path, "memory.swap.max")
	// 写入禁用swap配置
	if err := os.WriteFile(swapPath, []byte("0"), 0644); err != nil {
		// 记录错误日志
		logger.Error("Failed to disable swap: %v", err)
		// 返回错误
		return err
	}


	return nil
}

// SetCPULimit 设置CPU限制
// 参数：
// cpusStr：表示CPU限制的字符串，格式为 "cpu_quota cpu_period"
// 返回值：
// error：如果设置失败，则返回错误信息；否则返回nil
func (c *CGroupManager) SetCPULimit(cpusStr string) error {
	// 设置CPU限制
	// 解析传入的CPU字符串
	cpuQuota, cpuPeriod, err := ParseCPUs(cpusStr)
	if err != nil {
		logger.Error("Error parsing cpus: %v", err)
		return err
	}

	// 拼接CPU限制路径
	cpuMaxPath := filepath.Join(c.path, CpuMax)

	// 格式化CPU限制字符串
	cpuLimit := fmt.Sprintf("%d %d", cpuQuota, cpuPeriod)

	// 写入CPU限制到文件
	if err := os.WriteFile(cpuMaxPath, []byte(cpuLimit), 0644); err != nil {
		return err
	}

	return nil
}


// Cleanup 删除由c.path指定的目录及其所有子目录和文件。
// 如果删除过程中发生错误，会记录错误日志并返回错误。
// 如果没有错误发生，则返回nil。
func (c *CGroupManager) Cleanup() error {
	// 删除c.path指定的目录及其所有子目录和文件
	err := os.RemoveAll(c.path)
	if err != nil {
		// 记录错误日志
		logger.Error("Error cleaning up cgroup: %v", err)
		// 返回错误
		return err
	}
	return nil
}


// ParseCPUs 将 --cpus 字符串解析为 quota 和 period
// ParseCPUs 解析字符串形式的CPU值，并返回CPU配额（quota）和周期（period）
//
// 参数：
//   cpusStr: string类型，代表CPU值的字符串，支持浮点数输入，例如 "1.5"
//
// 返回值：
//   int: 返回计算出的CPU配额（quota）
//   int: 返回周期（period），固定为100ms（即100000微秒）
//   error: 返回解析过程中可能发生的错误，如果解析成功，则返回nil
func ParseCPUs(cpusStr string) (int, int, error) {
	// 支持浮点数输入，例如 "1.5"
	cpusFloat, err := strconv.ParseFloat(cpusStr, 64)
	if err != nil {
		// 解析错误，返回错误信息
		return 0, 0, fmt.Errorf("invalid cpus value: %s", cpusStr)
	}


	if cpusFloat <= 0 {
		// 如果输入的CPU值小于等于0，返回错误信息
		return 0, 0, fmt.Errorf("cpus must be greater than 0")
	}


	// 固定周期为 100ms（CGroups 推荐值）
	const period = 100000 // microseconds
	// 计算quota值
	quota := int(cpusFloat * float64(period))


	// 防止溢出或非法值
	if quota <= 0 {
		// 如果计算出的quota值小于等于0，返回错误信息
		return 0, 0, fmt.Errorf("calculated quota is invalid: %d", quota)
	}


	return quota, period, nil
}

