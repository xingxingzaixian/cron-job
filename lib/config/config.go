package config

import (
	"cronJob/internal/utils"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(configFile string) {
	// 设置默认配置
	setDefaultConfig()

	// 配置文件路径设置
	if configFile != "" {
		// 显式指定时基于基础目录解析相对路径
		viper.SetConfigFile(utils.BasePath(configFile))
		zap.S().Infof("使用指定配置文件: %s", configFile)
	} else {
		viper.AddConfigPath(utils.BaseDir())
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()
		zap.S().Infof("使用默认配置文件路径: %s", utils.BasePath("config.yaml"))
	}

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认配置继续运行
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if configFile != "" {
				zap.S().Warnf("指定的配置文件[%s]不存在，使用默认配置启动", configFile)
			} else {
				zap.S().Warn("配置文件不存在，使用默认配置启动")
			}
			return
		}
		// 其他错误（如权限问题、格式错误等）也使用默认配置继续运行
		zap.S().Errorf("读取配置文件失败: %v，使用默认配置继续运行", err)
		return
	}

	// 配置文件读取成功
	zap.S().Infof("配置文件读取成功: %s", viper.ConfigFileUsed())
}

// 设置默认配置值
func setDefaultConfig() {
	// HTTP 服务器默认配置
	viper.SetDefault("debug", "release")
	viper.SetDefault("http.addr", ":8210")
	viper.SetDefault("http.read_timeout", 10)
	viper.SetDefault("http.write_timeout", 10)
	viper.SetDefault("http.max_header_bytes", 20)

	// Swagger 默认配置
	viper.SetDefault("swagger.title", "定时任务服务swagger API")
	viper.SetDefault("swagger.desc", "这是一个简单的定时任务执行系统")
	viper.SetDefault("swagger.host", "127.0.0.1:8210")
	viper.SetDefault("swagger.base_path", "")

	// CORS 默认配置
	viper.SetDefault("cors.allowed_origins", "http://localhost:8210,http://127.0.0.1:8210")

	// 数据库默认配置（空值，表示未配置）
	viper.SetDefault("db.engine", "")
	viper.SetDefault("db.prefix", "sched_")
	// 数据目录基于基础目录，避免从其他目录启动时找不到数据库
	viper.SetDefault("db.data_dir", utils.BasePath("data"))

	// 认证默认配置
	viper.SetDefault("auth.enable", true)

	// 日志清理默认配置
	viper.SetDefault("log.retention_days", 30)               // 任务日志保留天数，0 表示不自动清理
	viper.SetDefault("log.cleanup_interval_hours", 6)        // 日志清理检查间隔（小时）
	viper.SetDefault("db.sqlite.vacuum_interval_hours", 168) // SQLite VACUUM 间隔（小时）

	// JWT 默认配置
	viper.SetDefault("jwt.secret", utils.GenerateRandomJWTSecret())
	viper.SetDefault("jwt.expires", 7200)
}
