package database

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySqlDB struct{}

func (m *MySqlDB) Create(option *gorm.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.name"),
	)
	zap.S().Info("mysql dsn: ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), option)
	if err != nil {
		zap.S().Fatalf("MySQL数据库连接失败:%s:%d", viper.GetString("db.host"), viper.GetInt("db.port"))
	}
	return db, err
}
