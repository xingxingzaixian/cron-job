package database

import (
	"cronJob/internal/global"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database interface {
	Create(option *gorm.Config) (*gorm.DB, error)
}

func InitDB(prefix string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})

	option := gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: true,
		},
	}

	// 使用工厂模式创建数据库实例
	factory := &DatabaseFactory{}

	// 验证配置
	if err := factory.ValidateConfig(); err != nil {
		panic(fmt.Sprintf("数据库配置验证失败: %v", err))
	}

	// 创建数据库实例
	database, err := factory.CreateDatabase()
	if err != nil {
		panic(fmt.Sprintf("创建数据库实例失败: %v", err))
	}

	// 显示连接信息
	zap.S().Infof("正在连接数据库: %s", factory.GetConnectionInfo())

	db, err := database.Create(&option)
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	// 配置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池参数
	maxIdleConns := viper.GetInt("db.max_idle_conns")
	maxOpenConns := viper.GetInt("db.max_open_conns")
	connMaxLifetime := viper.GetInt("db.conn_max_lifetime")

	if maxIdleConns <= 0 {
		maxIdleConns = 10
	}
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}
	if connMaxLifetime <= 0 {
		connMaxLifetime = 3600
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	zap.S().Infof("数据库连接池配置: MaxIdle=%d, MaxOpen=%d, MaxLifetime=%ds",
		maxIdleConns, maxOpenConns, connMaxLifetime)

	// 使用迁移管理器执行数据库迁移
	migrationManager := NewMigrationManager(db)

	// 检查数据库版本
	if err := migrationManager.CheckDatabaseVersion(); err != nil {
		zap.S().Warnf("检查数据库版本失败: %v", err)
	}

	// 执行自动迁移
	if err := migrationManager.AutoMigrate(); err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	// 创建初始数据
	if err := migrationManager.CreateInitialData(); err != nil {
		zap.S().Warnf("创建初始数据失败: %v", err)
	}

	global.GormDB = db
}
