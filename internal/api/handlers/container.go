package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phper95/tinydocker/container/models"
	"github.com/phper95/tinydocker/internal/api/errdefs"
	"github.com/phper95/tinydocker/internal/api/types"
	"net/http"
	"path/filepath"
)

// ListContainers 列出所有容器
func ListContainers(c *gin.Context) {
	containers := models.ReadContainersInfo()
	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, containers, nil))
}

func GetContainerInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, types.Error(errdefs.ErrInvalidContainerID, "容器id无效", "容器id不能为空"))
		return
	}
	// id 是容器ID，需要拼接配置文件路径
	filePath := filepath.Join(models.DefaultContainerInfoPath, id, models.DefaultContainerInfoFileName)
	info, err := models.ReadContainerInfo(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, types.Error(errdefs.ErrContainerNotFound, "容器不存在", err.Error()))
		return
	}
	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, info, nil))
}
