package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/phper95/tinydocker/internal/api/types"
	"net/http"
)

func ListImages(c *gin.Context) {
	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, "images", nil))
}
