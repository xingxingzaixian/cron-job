package cmd

import (
	"cronJob/internal/global"
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

	// 2. 检查是否需要安装
	needInstall := checkNeedInstall()
	global.InstallMode = needInstall
	global.IsInstalled = !needInstall
	
	// 3. 如果不需要安装，则初始化数据库和定时任务
	if !needInstall {
		// 初始化数据库
		database.InitDB(viper.GetString("db.prefix"))
		
		// 启动定时任务调度
		go func() {
			task.CronServerRun()
		}()
	}

	// 4. 启动web服务（无论是否需要安装都要启动）
	go func() {
		router.HttpServerRun(needInstall)
	}()

	// 5. 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Shutting down server...")
	if !needInstall {
		task.CronServerStop()
	}
	router.HttpServerStop()
}

// 检查是否需要安装
func checkNeedInstall() bool {
	// 检查数据库配置是否存在
	engine := viper.GetString("db.engine")
	if engine == "" {
		zap.S().Info("数据库未配置，进入安装模式")
		return true
	}
	
	switch engine {
	case "mysql", "postgresql", "postgres":
		host := viper.GetString("db.host")
		name := viper.GetString("db.name")
		user := viper.GetString("db.user")
		if host == "" || name == "" || user == "" {
			zap.S().Info("数据库配置不完整，进入安装模式")
			return true
		}
	case "sqlite", "sqlite3":
		name := viper.GetString("db.name")
		path := viper.GetString("db.path")
		if name == "" && path == "" {
			zap.S().Info("SQLite数据库配置不完整，进入安装模式")
			return true
		}
	default:
		zap.S().Info("不支持的数据库类型，进入安装模式")
		return true
	}
	
	zap.S().Info("数据库配置完整，正常启动")
	return false
}
