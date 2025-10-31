package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phper95/tinydocker/internal/api/errdefs"
	"github.com/phper95/tinydocker/internal/api/types"
	"github.com/phper95/tinydocker/internal/model"
	"github.com/phper95/tinydocker/internal/service"
	"github.com/phper95/tinydocker/pkg/logger"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将认证服务放入上下文，供后续权限中间件使用
		c.Set("auth_service", authService)
		// 获取认证信息
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("缺少认证信息")
			c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "缺少认证信息", ""))
			c.Abort()
			return
		}

		var authContext *model.AuthContext
		var err error

		// 解析认证头
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("不支持的认证方式")
			c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "不支持的认证方式", ""))
			c.Abort()
			return

		}
		// JWT Token 认证
		token := strings.TrimPrefix(authHeader, "Bearer ")
		authContext, err = authService.ValidateToken(token)
		if err != nil {
			logger.Error("认证失败: %v", err)
			c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrAuthFailed, "认证失败", err.Error()))
			c.Abort()
			return
		}

		// 将认证上下文存储到请求中
		c.Set("auth_context", authContext)
		c.Next()
	}
}

func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := c.Get("auth_context")
		logger.Debug("auth_context: %v", authContext)
		if !exists {
			logger.Error("缺少认证信息")
			c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "缺少认证信息", ""))
			c.Abort()
			return
		}

		ctx, ok := authContext.(*model.AuthContext)
		if !ok {
			logger.Error("认证上下文无效")
			c.JSON(http.StatusUnauthorized, types.Error(errdefs.ErrUnauthorized, "认证上下文无效", ""))
			c.Abort()
			return
		}

		// 获取认证服务
		authService, exists := c.Get("auth_service")
		if !exists {
			logger.Error("认证服务不可用")
			c.JSON(http.StatusInternalServerError, types.Error(errdefs.ErrAuthServiceUnavailable, "认证服务不可用", ""))
			c.Abort()
			return
		}

		s, ok := authService.(*service.AuthService)
		if !ok {
			logger.Error("认证服务类型错误")
			c.JSON(http.StatusInternalServerError, types.Error(errdefs.ErrAuthServiceTypeInvalid, "认证服务类型错误", ""))
			c.Abort()
			return
		}

		// 检查权限
		if !s.CheckPermission(ctx, resource, action) {
			logger.Error("权限不足: resource=%s, action=%s", resource, action)
			c.JSON(http.StatusForbidden, types.Error(errdefs.ErrAccessDenied, "权限不足", ctx.Username+" 无权访问 "+resource+" "+action))
			c.Abort()
			return
		}

		c.Next()
	}
}
