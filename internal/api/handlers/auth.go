package handlers

import (
	"github.com/phper95/tinydocker/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phper95/tinydocker/internal/api/errdefs"
	"github.com/phper95/tinydocker/internal/api/types"
	"github.com/phper95/tinydocker/internal/model"
	"github.com/phper95/tinydocker/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("登录请求参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, types.Error(errdefs.ErrInvalidParameter, "请求参数无效", err.Error()))
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		logger.Error("用户登录失败: %v", err)
		c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrLoginFailed, "登录失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, response, nil))
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("注册请求参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, types.Error(errdefs.ErrInvalidParameter, "请求参数无效", err.Error()))
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		logger.Error("用户注册失败: %v", err)
		c.JSON(http.StatusBadRequest, types.Error(errdefs.ErrRegisterFailed, "注册失败", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, types.Success(types.ApiVersionV1, user, nil))
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	authContext, exists := c.Get("auth_context")
	if !exists {
		logger.Error("缺少认证信息")
		c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "缺少认证信息", ""))
		return
	}

	ctx, ok := authContext.(*model.AuthContext)
	if !ok {
		logger.Error("认证上下文无效")
		c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "认证上下文无效", ""))
		return
	}

	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, ctx, nil))
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从 Authorization 头中提取 Bearer Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		logger.Error("缺少或无效的认证信息")
		c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "缺少或无效的认证信息", ""))
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	err := h.service.Logout(token)
	if err != nil {
		logger.Error("用户登出失败: %v", err)
		c.JSON(http.StatusInternalServerError, types.Error(errdefs.ErrLogoutFailed, "登出失败", err.Error()))
		return
	}

	c.JSON(http.StatusOK, types.Success(types.ApiVersionV1, nil, nil))

}
