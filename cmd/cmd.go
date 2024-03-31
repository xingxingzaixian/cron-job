package cmd

import (
	task "cronJob/internal/service/cron/task_manager"
	"cronJob/internal/service/web/router"
	"cronJob/lib/config"
	"cronJob/lib/database"
	"cronJob/lib/logger"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var AppVersion = "1.0.0"

var rootCmd = &cobra.Command{
	Use:     "CronJob",
	Short:   "CronJob定时任务管理系统",
	Long:    `CronJob是用Go实现的秒级分布式定时任务执行管理系统`,
	Version: AppVersion,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.PersistentFlags().GetString("config")
		startServer(configFile)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "config file path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startServer(configFile string) {
	// 0. 初始化日志
	logger.InitLogger()

	// 1. 初始化配置
	config.InitConfig(configFile)

	// 2. 初始化数据库
	database.InitDB(viper.GetString("db.prefix"))

	// 3. 启动配置web服务
	go func() {
		router.HttpServerRun()
	}()

	// 4. 启动定时任务调度
	go func() {
		task.CronServerRun()
	}()

	// 5. 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Shutting down server...")
	task.CronServerStop()
	router.HttpServerStop()
}
