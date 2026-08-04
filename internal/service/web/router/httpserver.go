package router

import (
	"context"
	"cronJob/internal/global"
	middleware2 "cronJob/internal/service/web/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func HttpServerRun(installMode bool) {
	gin.SetMode(viper.GetString("debug"))
	router := InitRouter(installMode, middleware2.Cors(), middleware2.TranslationMiddleware())

	// 设置默认配置值，防止配置文件不完整时出错
	addr := viper.GetString("http.addr")
	if addr == "" {
		addr = ":8210"
	}

	readTimeout := viper.GetInt("http.read_timeout")
	if readTimeout <= 0 {
		readTimeout = 10
	}

	writeTimeout := viper.GetInt("http.write_timeout")
	if writeTimeout <= 0 {
		writeTimeout = 10
	}

	maxHeaderBytes := viper.GetUint("http.max_header_bytes")
	if maxHeaderBytes == 0 {
		maxHeaderBytes = 20
	}

	global.HttpSrvHandler = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << maxHeaderBytes,
	}

	go func() {
		if installMode {
			zap.S().Infof(" [INFO] HttpServerRun (Install Mode):%s\n", addr)
		} else {
			zap.S().Infof(" [INFO] HttpServerRun:%s\n", addr)
		}
		if err := global.HttpSrvHandler.ListenAndServe(); err != nil {
			zap.S().Errorf(" [ERROR] HttpServerRun:%s err:%v\n", addr, err)
		}
	}()
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := global.HttpSrvHandler.Shutdown(ctx); err != nil {
		zap.S().Errorf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	zap.S().Infof(" [INFO] HttpServerStop stopped\n")
}
