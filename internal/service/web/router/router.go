package router

import (
	"cronJob/docs"
	"cronJob/internal/service/web/api"
	middleware2 "cronJob/internal/service/web/middleware"
	"cronJob/web"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/fs"
	"net/http"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = viper.GetString("swagger.title")
	docs.SwaggerInfo.Description = viper.GetString("swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = viper.GetString("swagger.host")
	docs.SwaggerInfo.BasePath = viper.GetString("swagger.base")
	docs.SwaggerInfo.Schemes = []string{"http"}

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery(), gzip.Gzip(gzip.DefaultCompression))
	router.Use(middlewares...)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/admin"
		router.HandleContext(ctx)
	})

	webFs, _ := fs.Sub(web.WebFS, "dist")
	router.StaticFS("/admin", http.FS(webFs))

	apiRouter := router.Group("/api")
	loginRouter := apiRouter.Group("/")
	{
		api.LoginRegister(loginRouter)
	}

	apiRouter.Use(middleware2.JwtAuthMiddleware())
	taskRouter := apiRouter.Group("/task")
	{
		api.TaskRegister(taskRouter)
	}

	taskLogRouter := apiRouter.Group("/taskLog")
	{
		api.TaskLogRegister(taskLogRouter)
	}

	userRouter := apiRouter.Group("/user")
	{
		api.UserRegister(userRouter)
	}

	roleRouter := apiRouter.Group("/role")
	{
		api.RoleRegister(roleRouter)
	}
	return router
}
