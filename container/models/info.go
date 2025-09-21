package models

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/phper95/tinydocker/pkg/logger"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"
)

const (
	DefaultContainerInfoPath     = "/var/lib/tinydocker/containers"
	DefaultImagePath             = "/var/lib/tinydocker/image"
	DefaultContainerInfoFileName = "config.json"
	ContainerStateRunning        = "running"
	ContainerStateStopped        = "stopped"
)

type Info struct {
	Id          string   `json:"id"`          // 容器Id
	Name        string   `json:"name"`        // 容器名
	Pid         int      `json:"pid"`         // 容器的init进程在宿主机上的 PID
	Command     string   `json:"command"`     // 容器内init运行命令
	State       string   `json:"state"`       // 容器的状态
	StartedAt   string   `json:"started_at"`  // 启动时间
	FinishedAt  string   `json:"finished_at"` // 结束时间
	Image       string   `json:"image"`       // 容器使用的镜像名称
	Network     string   `json:"network"`
	IpAddress   string   `json:"ipAddress"`
	PortMapping []string `json:"port_mapping"` // 端口映射
}

func WriteContainerInfo(info *Info) error {
	// 序列化Info结构体到json字符串
	jsonStr, err := json.Marshal(info)
	if err != nil {
		return err
	}
	// 写入json字符串到文件
	dirPath := filepath.Join(DefaultContainerInfoPath, info.Id)
	if err := os.MkdirAll(dirPath, 0622); err != nil {
		logger.Error("mkdirall error: ", err)
		return err
	}
	filePath := filepath.Join(dirPath, DefaultContainerInfoFileName)
	err = os.WriteFile(filePath, jsonStr, 0622)
	if err != nil {
		logger.Error("writefile error: ", err)
		return err
	}
	return nil
}

func UpdateContainerState(containerID string, state string) error {
	// 构建容器配置文件路径
	filePath := filepath.Join(DefaultContainerInfoPath, containerID, DefaultContainerInfoFileName)

	// 读取现有配置
	info, err := ReadContainerInfo(filePath)
	if err != nil {
		return err
	}

	// 更新状态
	info.State = state
	if state == ContainerStateStopped {
		info.FinishedAt = time.Now().Format(time.DateTime)
	}

	// 写回文件
	return WriteContainerInfo(info)
}

func GenerateRandomContainerID() string {
	bytes := make([]byte, 32) // 64个十六进制字符
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", bytes)
}

func PrintContainersInfo() error {
	containersInfo := ReadContainersInfo()
	if len(containersInfo) == 0 {
		return nil
	}
	// 格式化输出表格
	tableWri := tabwriter.NewWriter(os.Stdout, 6, 2, 1, '\t', 0)
	fmt.Fprintln(tableWri, "ID\tNAME\tPID\tCOMMAND\tSTATE\tSTARTED_AT\tFINISHED_AT")
	for _, info := range containersInfo {
		fmt.Fprintf(tableWri, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			info.Id, info.Name, info.Pid, info.Command, info.State, info.StartedAt, info.FinishedAt)
	}
	if err := tableWri.Flush(); err != nil {
		logger.Error("flush error: ", err)
		return err
	}
	return nil
}

func ReadContainersInfo() []Info {
	// 读取容器信息目录
	dirs, err := os.ReadDir(DefaultContainerInfoPath)
	if err != nil {
		logger.Error("readdir error: ", err)
		return nil
	}
	// 遍历目录，读取每个容器的配置信息
	var infos []Info
	for _, d := range dirs {
		logger.Debug("dir name: ", d.Name())
		if !d.IsDir() {
			continue
		}
		filePath := filepath.Join(DefaultContainerInfoPath, d.Name(), DefaultContainerInfoFileName)
		info, err := ReadContainerInfo(filePath)
		if err != nil {
			logger.Error("readcontainerinfo error: ", err)
			continue
		}
		infos = append(infos, *info)
	}
	return infos
}

func ReadContainerInfo(filePath string) (*Info, error) {
	// 读取json字符串
	jsonStr, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("readfile error: ", err)
		return nil, err
	}
	// 反序列化json字符串到Info结构体
	var info Info
	if err := json.Unmarshal(jsonStr, &info); err != nil {
		logger.Error("unmarshal error: ", err)
		return nil, err
	}
	return &info, nil
}
