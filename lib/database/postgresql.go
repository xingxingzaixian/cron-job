package database

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDB struct{}

func (p *PostgreSQLDB) Create(option *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.name"),
		viper.GetInt("db.port"),
	)
	
	zap.S().Info("postgresql dsn: ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), option)
	if err != nil {
		zap.S().Fatalf("PostgreSQL数据库连接失败:%s:%d", viper.GetString("db.host"), viper.GetInt("db.port"))
	}
	return db, err
}
