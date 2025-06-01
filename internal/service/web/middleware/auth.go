package middleware

import (
	jwt2 "cronJob/lib/jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

// AuthMiddleware 认证中间件，根据配置决定是否启用认证
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查是否启用认证
		if !viper.GetBool("auth.enable") {
			// 如果未启用认证，直接跳过认证检查
			ctx.Next()
			return
		}

		// 启用认证时，执行JWT认证逻辑
		authHeader := ctx.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[1] == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		claims, err := jwt2.ParseToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}

// IsAuthEnabled 检查是否启用了认证功能
func IsAuthEnabled() bool {
	return viper.GetBool("auth.enable")
}
