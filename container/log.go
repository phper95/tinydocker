package container

import (
	"bufio"
	"fmt"
	"github.com/phper95/tinydocker/container/models"
	"github.com/phper95/tinydocker/pkg/logger"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultContainerLogFileName = "container.log"
)

func PrintContainerLogs(containerID string, follow bool) error {
	logDir := filepath.Join(models.DefaultContainerInfoPath, containerID)
	logFilePath := filepath.Join(logDir, DefaultContainerLogFileName)
	// Check if log file exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		logger.Error("log file for container %s does not exist", containerID)
		return fmt.Errorf("log file for container %s does not exist", containerID)
	}
	// Open log file
	file, err := os.Open(logFilePath)
	if err != nil {
		logger.Error("failed to open log file: %v", err)
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()
	// If follow flag is set, continuously read new logs
	if follow {
		// Move to the end of the file
		_, err = file.Seek(0, io.SeekEnd)
		if err != nil {
			return fmt.Errorf("failed to seek to end of file: %v", err)
		}

		// Continuously read new content
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(1 * time.Second) // Wait before retrying
					continue
				} else {
					return fmt.Errorf("error reading log file: %v", err)
				}
			}
			fmt.Print(line)
		}
	} else {
		// Read and print entire log file
		content, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read log file: %v", err)
		}
		fmt.Print(string(content))
	}

	return nil
}
