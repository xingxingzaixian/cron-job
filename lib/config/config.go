package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(configFile string) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()
	}

	// 设置默认配置值
	setDefaultConfig()

	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认配置继续运行
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zap.S().Warnf("配置文件[%s]不存在，使用默认配置启动", configFile)
			return
		}
		// 其他错误则退出
		zap.S().Panicf("读取配置文件[%s]失败: %v", configFile, err)
	}

	zap.S().Infof("配置文件【%s】读取成功", viper.ConfigFileUsed())
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
	viper.SetDefault("cors.allowed_origins", "http://localhost:3000,http://127.0.0.1:3000,http://localhost:8210,http://127.0.0.1:8210")
	
	// 数据库默认配置（空值，表示未配置）
	viper.SetDefault("db.engine", "")
	viper.SetDefault("db.prefix", "sched_")
	
	// 认证默认配置
	viper.SetDefault("auth.enable", true)
	
	// JWT 默认配置
	viper.SetDefault("jwt.secret", "default_jwt_secret_change_me")
	viper.SetDefault("jwt.expires", 7200)
}
