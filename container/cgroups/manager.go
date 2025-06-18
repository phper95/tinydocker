package cgroups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type CGroupManager struct {
	path string
}

func NewCGroupManager(name string) *CGroupManager {
	path := filepath.Join("/sys/fs/cgroup", name)
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}
	return &CGroupManager{path: path}
}

func (c *CGroupManager) Apply(pid int) {
	pidStr := strconv.Itoa(pid)
	if err := ioutil.WriteFile(
		filepath.Join(c.path, "cgroup.procs"),
		[]byte(pidStr),
		0644,
	); err != nil {
		panic(err)
	}
}

func (c *CGroupManager) SetCPULimit(percent int) {
	limit := fmt.Sprintf("%d", percent)
	if err := ioutil.WriteFile(
		filepath.Join(c.path, "cpu.max"),
		[]byte(limit+" 100000"),
		0644,
	); err != nil {
		panic(err)
	}
}

func (c *CGroupManager) Cleanup() {
	os.RemoveAll(c.path)
}
