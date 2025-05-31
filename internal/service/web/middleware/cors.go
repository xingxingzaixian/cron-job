package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")

		// 从配置文件读取允许的域名，如果没有配置则使用默认值
		allowedOrigins := viper.GetString("cors.allowed_origins")
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		// 检查是否允许该域名
		if allowedOrigins != "*" {
			origins := strings.Split(allowedOrigins, ",")
			allowed := false
			for _, allowedOrigin := range origins {
				if strings.TrimSpace(allowedOrigin) == origin {
					allowed = true
					break
				}
			}
			if allowed {
				ctx.Header("Access-Control-Allow-Origin", origin)
			}
		} else {
			ctx.Header("Access-Control-Allow-Origin", "*")
		}

		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Max-Age", "86400")

		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}
