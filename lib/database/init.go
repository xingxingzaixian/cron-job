package database

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
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

	db.AutoMigrate(&models.Task{}, &models.TaskLog{}, &models.User{})
	global.GormDB = db
}
