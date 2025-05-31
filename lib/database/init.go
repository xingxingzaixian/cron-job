package database

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database interface {
	Create(option *gorm.Config) *gorm.DB
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

	sql := MySqlDB{}
	db, err := sql.Create(&option)
	if err != nil {
		panic(err)
	}

	// 配置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(viper.GetInt("db.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("db.max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("db.conn_max_lifetime")) * time.Second)

	db.AutoMigrate(&models.Task{}, &models.TaskLog{}, &models.User{}, &models.Role{})
	global.GormDB = db
}
