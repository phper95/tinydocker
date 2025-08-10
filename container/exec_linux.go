//go:build linux && cgo

package container

/*
#define _GNU_SOURCE
#include <sched.h>
#include <unistd.h>
#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <stdlib.h>

static int setns_cgo(const char* nspath, int nstype) {
    int fd = open(nspath, O_RDONLY);
    if (fd < 0) {
        return -1;
    }
    int result = setns(fd, nstype);
    int errno_copy = errno;
    close(fd);
    errno = errno_copy;
    return result;
}
*/
import "C"

import (
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"unsafe"
)

// ExecContainer is the internal entrypoint to join target namespaces and run the user command.
// It expects env TINYDOCKER_EXEC_PID to be set and args to contain the command to run.
func ExecContainer(args cli.Args) error {
	pidStr := os.Getenv(ExecTargetPidEnv)
	if pidStr == "" {
		return fmt.Errorf("%s not set", ExecTargetPidEnv)
	}
	targetPid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("invalid %s: %v", ExecTargetPidEnv, err)
	}
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	// join namespaces of target pid (mnt, uts, ipc, net). pid ns join affects children only.
	namespaces := []string{"mnt", "uts", "ipc", "net"}
	for _, ns := range namespaces {
		if err := joinNamespace(targetPid, ns); err != nil {
			logger.Error("joinNamespace(%s) failed: %v", ns, err)
			return fmt.Errorf("join %s namespace failed: %w", ns, err)
		}
	}

	// chdir to root in the target mount namespace to mimic container context
	_ = os.Chdir("/")

	// Find and exec the command in-place
	cmdPath, err := lookupPath(args.Get(0))
	if err != nil {
		return err
	}

	argv := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		argv = append(argv, args.Get(i))
	}

	// replace current process image
	return syscall.Exec(cmdPath, argv, os.Environ())
}

func joinNamespace(pid int, ns string) error {
	// 检查目标进程是否存在
	procPath := filepath.Join("/proc", strconv.Itoa(pid))
	if _, err := os.Stat(procPath); err != nil {
		logger.Error("target process %d does not exist: %v", pid, err)
		return fmt.Errorf("target process %d does not exist: %w", pid, err)
	}

	nsPath := filepath.Join("/proc", strconv.Itoa(pid), "ns", ns)

	// 检查命名空间文件是否存在
	if _, err := os.Stat(nsPath); err != nil {
		return fmt.Errorf("namespace file %s does not exist or cannot be accessed: %w", nsPath, err)
	}

	// 添加调试信息
	logger.Debug("Attempting to join namespace: %s", nsPath)

	// 根据命名空间类型设置nstype参数
	var nstype int
	switch ns {
	case "mnt":
		nstype = syscall.CLONE_NEWNS
	case "uts":
		nstype = syscall.CLONE_NEWUTS
	case "ipc":
		nstype = syscall.CLONE_NEWIPC
	case "net":
		nstype = syscall.CLONE_NEWNET
	case "pid":
		nstype = syscall.CLONE_NEWPID
	case "user":
		nstype = syscall.CLONE_NEWUSER
	case "cgroup":
		nstype = syscall.CLONE_NEWCGROUP
	default:
		nstype = 0 // 让内核自动检测
	}

	logger.Debug("Using nstype value: %d for namespace: %s", nstype, ns)

	// 将Go字符串转换为C字符串
	cNsPath := C.CString(nsPath)
	defer C.free(unsafe.Pointer(cNsPath))

	// 使用C代码处理整个操作，包括打开文件和调用setns
	if ret := C.setns_cgo(cNsPath, C.int(nstype)); ret != 0 {
		logger.Error("setns(%s) failed with nstype=%d ret=%v", nsPath, nstype, syscall.Errno(ret))
		// 尝试使用nstype=0
		// if ret2 := C.setns_cgo(cNsPath, C.int(0)); ret2 != 0 {
		// 	logger.Error("setns(%s) failed with nstype=0: %v", nsPath, syscall.Errno(ret2))
		// 	return fmt.Errorf("setns(%s) failed: %w", nsPath, syscall.Errno(ret2))
		// }
	}
	logger.Debug("joined namespace: ", nsPath)
	return nil
}

func lookupPath(file string) (string, error) {
	if filepath.IsAbs(file) {
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}
	// simple PATH search using current env
	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		p := filepath.Join(dir, file)
		if st, err := os.Stat(p); err == nil && !st.IsDir() && (st.Mode()&0111) != 0 {
			return p, nil
		}
	}
	return "", fmt.Errorf("executable %s not found in PATH", file)
}
