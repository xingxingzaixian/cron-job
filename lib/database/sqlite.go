package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // 纯Go SQLite驱动
)

// SQLiteDB SQLite数据库实现
type SQLiteDB struct{}

// Create 创建SQLite数据库连接
func (s *SQLiteDB) Create(option *gorm.Config) (*gorm.DB, error) {
	// 获取数据库文件路径
	dbPath := s.getDBPath()

	// 确保数据库目录存在
	if err := s.ensureDBDir(dbPath); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %v", err)
	}

	// 记录连接信息（不包含敏感信息）
	zap.S().Infof("SQLite 数据库文件路径: %s", dbPath)

	// 构建连接字符串，指定使用 modernc.org/sqlite 驱动
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", dbPath)

	// 打开数据库连接，使用纯Go驱动
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dsn,
	}, option)
	if err != nil {
		zap.S().Errorf("SQLite数据库连接失败: %s, 错误: %v", dbPath, err)
		return nil, fmt.Errorf("SQLite数据库连接失败: %v", err)
	}

	// 配置SQLite特定的设置
	if err := s.configureSQLite(db); err != nil {
		zap.S().Errorf("配置SQLite失败: %v", err)
		return nil, fmt.Errorf("配置SQLite失败: %v", err)
	}

	zap.S().Info("SQLite数据库连接成功")
	return db, nil
}

// getDBPath 获取数据库文件路径
func (s *SQLiteDB) getDBPath() string {
	// 优先使用配置文件中的路径
	if dbPath := viper.GetString("db.path"); dbPath != "" {
		return dbPath
	}

	// 使用数据库名称作为文件名
	dbName := viper.GetString("db.name")
	if dbName == "" {
		dbName = "cronJob"
	}

	// 如果没有扩展名，添加.db扩展名
	if filepath.Ext(dbName) == "" {
		dbName += ".db"
	}

	// 默认存储在data目录下
	dataDir := viper.GetString("db.data_dir")
	if dataDir == "" {
		dataDir = "data"
	}

	return filepath.Join(dataDir, dbName)
}

// ensureDBDir 确保数据库目录存在
func (s *SQLiteDB) ensureDBDir(dbPath string) error {
	dir := filepath.Dir(dbPath)

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}
		zap.S().Infof("创建数据库目录: %s", dir)
	}

	return nil
}

// configureSQLite 配置SQLite特定的设置
func (s *SQLiteDB) configureSQLite(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 启用外键约束
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		zap.S().Warnf("启用外键约束失败: %v", err)
	}

	// 设置WAL模式以提高并发性能
	if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		zap.S().Warnf("设置WAL模式失败: %v", err)
	}

	// 设置同步模式
	syncMode := viper.GetString("db.sqlite.synchronous")
	if syncMode == "" {
		syncMode = "NORMAL" // 默认为NORMAL，平衡性能和安全性
	}
	if _, err := sqlDB.Exec(fmt.Sprintf("PRAGMA synchronous = %s", syncMode)); err != nil {
		zap.S().Warnf("设置同步模式失败: %v", err)
	}

	// 设置缓存大小
	cacheSize := viper.GetInt("db.sqlite.cache_size")
	if cacheSize <= 0 {
		cacheSize = -64000 // 默认64MB缓存
	}
	if _, err := sqlDB.Exec(fmt.Sprintf("PRAGMA cache_size = %d", cacheSize)); err != nil {
		zap.S().Warnf("设置缓存大小失败: %v", err)
	}

	// 设置临时存储
	tempStore := viper.GetString("db.sqlite.temp_store")
	if tempStore == "" {
		tempStore = "MEMORY" // 默认使用内存存储临时数据
	}
	if _, err := sqlDB.Exec(fmt.Sprintf("PRAGMA temp_store = %s", tempStore)); err != nil {
		zap.S().Warnf("设置临时存储失败: %v", err)
	}

	// 设置忙等待超时
	busyTimeout := viper.GetInt("db.sqlite.busy_timeout")
	if busyTimeout <= 0 {
		busyTimeout = 30000 // 默认30秒
	}
	if _, err := sqlDB.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d", busyTimeout)); err != nil {
		zap.S().Warnf("设置忙等待超时失败: %v", err)
	}

	zap.S().Info("SQLite配置完成")
	return nil
}

// GetDBInfo 获取数据库信息
func (s *SQLiteDB) GetDBInfo() map[string]interface{} {
	return map[string]interface{}{
		"type":     "sqlite",
		"path":     s.getDBPath(),
		"data_dir": viper.GetString("db.data_dir"),
		"name":     viper.GetString("db.name"),
	}
}

// CheckDBFile 检查数据库文件是否存在
func (s *SQLiteDB) CheckDBFile() bool {
	dbPath := s.getDBPath()
	_, err := os.Stat(dbPath)
	return !os.IsNotExist(err)
}

// GetDBSize 获取数据库文件大小
func (s *SQLiteDB) GetDBSize() (int64, error) {
	dbPath := s.getDBPath()
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// BackupDB 备份数据库文件
func (s *SQLiteDB) BackupDB(backupPath string) error {
	dbPath := s.getDBPath()

	// 确保备份目录存在
	if err := s.ensureDBDir(backupPath); err != nil {
		return err
	}

	// 读取源文件
	sourceData, err := os.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("读取数据库文件失败: %v", err)
	}

	// 写入备份文件
	if err := os.WriteFile(backupPath, sourceData, 0644); err != nil {
		return fmt.Errorf("写入备份文件失败: %v", err)
	}

	zap.S().Infof("数据库备份完成: %s -> %s", dbPath, backupPath)
	return nil
}

// RestoreDB 从备份恢复数据库
func (s *SQLiteDB) RestoreDB(backupPath string) error {
	dbPath := s.getDBPath()

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupPath)
	}

	// 读取备份文件
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %v", err)
	}

	// 确保数据库目录存在
	if err := s.ensureDBDir(dbPath); err != nil {
		return err
	}

	// 写入数据库文件
	if err := os.WriteFile(dbPath, backupData, 0644); err != nil {
		return fmt.Errorf("写入数据库文件失败: %v", err)
	}

	zap.S().Infof("数据库恢复完成: %s -> %s", backupPath, dbPath)
	return nil
}
