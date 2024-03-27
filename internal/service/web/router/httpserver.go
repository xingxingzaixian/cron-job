package router

import (
	"context"
	middleware2 "cronJob/internal/service/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var HttpSrvHandler *http.Server

func HttpServerRun() {
	gin.SetMode(viper.GetString("debug"))
	router := InitRouter(middleware2.Cors(), middleware2.TranslationMiddleware())
	HttpSrvHandler = &http.Server{
		Addr:           viper.GetString("http.addr"),
		Handler:        router,
		ReadTimeout:    viper.GetDuration("http.readtimeout") * time.Second,
		WriteTimeout:   viper.GetDuration("http.writetimeout") * time.Second,
		MaxHeaderBytes: 1 << viper.GetUint("http.maxheaderbytes"),
	}

	go func() {
		zap.S().Infof(" [INFO] HttpServerRun:%s\n", viper.GetString("http.addr"))
		if err := HttpSrvHandler.ListenAndServe(); err != nil {
			zap.S().Fatalf(" [ERROR] HttpServerRun:%s err:%v\n", viper.GetString("http.addr"), err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		zap.S().Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	zap.S().Infof(" [INFO] HttpServerStop stopped\n")
}
