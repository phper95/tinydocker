package cgroups

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
	"path/filepath"
	"strconv"
)

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
func (c *CGroupManager) SetCPULimit(cpuQuota, cpuPeriod int) error {
	// 设置CPU限制
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
