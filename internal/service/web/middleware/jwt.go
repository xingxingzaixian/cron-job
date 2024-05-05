package middleware

import (
	jwt2 "cronJob/lib/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
			fmt.Println(22222, parts)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		claims, err := jwt2.ParseToken(parts[1])
		if err != nil {
			fmt.Println(33333, parts[1], err.Error())
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
