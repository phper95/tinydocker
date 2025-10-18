package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phper95/tinydocker/container/models"
	"github.com/phper95/tinydocker/internal/api/types"
	"net/http"
)

// ListContainers 列出所有容器
func ListContainers(c *gin.Context) {
	containers := []models.Info{
		{
			Id:    "abc123",
			Name:  "test-container",
			State: "running",
			Image: "busybox:latest",
		},
	}
	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, containers, nil))
}
