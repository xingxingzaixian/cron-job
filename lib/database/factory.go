package database

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// DatabaseFactory 数据库工厂
type DatabaseFactory struct{}

// CreateDatabase 根据配置创建数据库实例
func (f *DatabaseFactory) CreateDatabase() (Database, error) {
	engine := viper.GetString("db.engine")
	
	switch engine {
	case "mysql":
		zap.S().Info("使用 MySQL 数据库")
		return &MySqlDB{}, nil
	case "postgresql", "postgres":
		zap.S().Info("使用 PostgreSQL 数据库")
		return &PostgreSQLDB{}, nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s，支持的类型: mysql, postgresql", engine)
	}
}

// GetSupportedEngines 获取支持的数据库引擎列表
func (f *DatabaseFactory) GetSupportedEngines() []string {
	return []string{"mysql", "postgresql", "postgres"}
}

// ValidateConfig 验证数据库配置
func (f *DatabaseFactory) ValidateConfig() error {
	engine := viper.GetString("db.engine")
	if engine == "" {
		return fmt.Errorf("数据库引擎配置不能为空")
	}
	
	host := viper.GetString("db.host")
	if host == "" {
		return fmt.Errorf("数据库主机配置不能为空")
	}
	
	port := viper.GetInt("db.port")
	if port <= 0 {
		return fmt.Errorf("数据库端口配置无效: %d", port)
	}
	
	name := viper.GetString("db.name")
	if name == "" {
		return fmt.Errorf("数据库名称配置不能为空")
	}
	
	user := viper.GetString("db.user")
	if user == "" {
		return fmt.Errorf("数据库用户名配置不能为空")
	}
	
	// 检查引擎是否支持
	supported := f.GetSupportedEngines()
	for _, supportedEngine := range supported {
		if engine == supportedEngine {
			return nil
		}
	}
	
	return fmt.Errorf("不支持的数据库引擎: %s，支持的引擎: %v", engine, supported)
}

// GetDefaultPort 根据数据库类型获取默认端口
func (f *DatabaseFactory) GetDefaultPort(engine string) int {
	switch engine {
	case "mysql":
		return 3306
	case "postgresql", "postgres":
		return 5432
	default:
		return 0
	}
}

// GetConnectionInfo 获取连接信息（用于日志显示，不包含密码）
func (f *DatabaseFactory) GetConnectionInfo() string {
	engine := viper.GetString("db.engine")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	name := viper.GetString("db.name")
	user := viper.GetString("db.user")
	
	return fmt.Sprintf("%s://%s@%s:%d/%s", engine, user, host, port, name)
}
