package container

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/phper95/tinydocker/pkg/logger"
)

const ExecTargetPidEnv = "TINYDOCKER_EXEC_PID"

// Exec runs a command inside a running container identified by name.
// It re-execs the current binary with the hidden command "exec-container",
// passing the target container's init PID via env for the child to join namespaces.
func Exec(name string, args []string, enableTTY bool) error {
	info, err := getContainerInfoByName(name)
	if err != nil {
		logger.Error("get container info failed: %v", err)
		return err
	}
	if info.Pid <= 0 {
		logger.Error("container %s is not running (no pid)", name)
		return fmt.Errorf("container %s is not running (no pid)", name)
	}
	logger.Debug("exec target pid:", info.Pid, "args:", args)

	// prepare re-exec
	cmd := exec.Command("/proc/self/exe", "exec-container")
	cmd.Args = append(cmd.Args, args...)

	// attach stdio (we don't allocate a pty here)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// pass env with target pid
	env := os.Environ()
	env = append(env, ExecTargetPidEnv+"="+strconv.Itoa(info.Pid))
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		logger.Error("exec in container %s failed: %v", name, err)
		return fmt.Errorf("exec in container %s failed: %w", name, err)
	}
	return nil
}

func getContainerInfoByName(name string) (*Info, error) {
	all := ReadContainersInfo()
	for _, c := range all {
		if strings.EqualFold(c.Name, name) {
			cc := c
			return &cc, nil
		}
	}
	return nil, fmt.Errorf("container with name %s not found", name)
}
