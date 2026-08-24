package cmd

import (
	"cronJob/internal/global"
	"cronJob/internal/service/cron/logcleaner"
	task "cronJob/internal/service/cron/task_manager"
	"cronJob/internal/service/web/router"
	"cronJob/internal/utils"
	"cronJob/lib/config"
	"cronJob/lib/database"
	"cronJob/lib/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	// 默认空值，由 InitConfig 基于基础目录解析
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path")
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

		// 迁移历史明文SSH密码为加密存储
		if n, err := task.MigrateSSHSecrets(); err != nil {
			zap.S().Warnf("SSH密码加密迁移失败: %v", err)
		} else if n > 0 {
			zap.S().Infof("SSH密码加密迁移完成，共 %d 条", n)
		}

		// 启动定时任务调度
		go func() {
			task.CronServerRun()
		}()

		// 启动任务日志自动清理
		logcleaner.Start()
	}

	// 4. 启动web服务（无论是否需要安装都要启动）
	go func() {
		router.HttpServerRun(needInstall)
	}()

	// 5. 监听退出信号
	// 使用带缓冲的通道，避免信号在 Notify 尚未读取时被丢弃
	quit := make(chan os.Signal, 1)
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
	// 检查install.lock文件是否存在
	_, err := os.Stat(utils.BasePath("install.lock"))
	if os.IsNotExist(err) {
		zap.S().Info("install.lock文件不存在，进入安装模式")
		return true
	}

	zap.S().Info("install.lock文件存在，系统已安装，正常启动")
	return false
}
