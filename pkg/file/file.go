package file

import (
	"github.com/phper95/tinydocker/pkg/logger"
	"os"
)

// 判断文件夹是否存在
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		logger.Error("failed to access path:", path, "error:", err)
		return false
	}
	return s.IsDir()
}
