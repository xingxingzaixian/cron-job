package database

import (
	"cronJob/internal/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrationManager 数据库迁移管理器
type MigrationManager struct {
	db     *gorm.DB
	engine string
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:     db,
		engine: viper.GetString("db.engine"),
	}
}

// AutoMigrate 执行自动迁移
func (m *MigrationManager) AutoMigrate() error {
	zap.S().Info("开始执行数据库迁移...")
	
	// 执行基本的自动迁移
	err := m.db.AutoMigrate(
		&models.Task{},
		&models.TaskLog{},
		&models.User{},
		&models.Role{},
	)
	if err != nil {
		return err
	}
	
	// 根据数据库类型执行特定的优化
	switch m.engine {
	case "mysql":
		err = m.optimizeForMySQL()
	case "postgresql", "postgres":
		err = m.optimizeForPostgreSQL()
	}
	
	if err != nil {
		zap.S().Warnf("数据库优化失败: %v", err)
		// 优化失败不影响主流程，只记录警告
	}
	
	zap.S().Info("数据库迁移完成")
	return nil
}

// optimizeForMySQL MySQL特定优化
func (m *MigrationManager) optimizeForMySQL() error {
	zap.S().Info("执行 MySQL 特定优化...")
	
	// MySQL特定的索引和优化
	sqls := []string{
		// 为任务表添加复合索引
		"CREATE INDEX IF NOT EXISTS idx_task_status_protocol ON " + m.getTableName("task") + " (status, protocol)",
		
		// 为任务日志表添加复合索引
		"CREATE INDEX IF NOT EXISTS idx_tasklog_task_status ON " + m.getTableName("task_log") + " (task_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_tasklog_start_time ON " + m.getTableName("task_log") + " (start_time)",
		
		// 设置表引擎和字符集
		"ALTER TABLE " + m.getTableName("task") + " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		"ALTER TABLE " + m.getTableName("task_log") + " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		"ALTER TABLE " + m.getTableName("user") + " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		"ALTER TABLE " + m.getTableName("role") + " ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
	}
	
	return m.executeSQLs(sqls)
}

// optimizeForPostgreSQL PostgreSQL特定优化
func (m *MigrationManager) optimizeForPostgreSQL() error {
	zap.S().Info("执行 PostgreSQL 特定优化...")
	
	// PostgreSQL特定的索引和优化
	sqls := []string{
		// 为任务表添加复合索引
		"CREATE INDEX IF NOT EXISTS idx_task_status_protocol ON " + m.getTableName("task") + " (status, protocol)",
		
		// 为任务日志表添加复合索引
		"CREATE INDEX IF NOT EXISTS idx_tasklog_task_status ON " + m.getTableName("task_log") + " (task_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_tasklog_start_time ON " + m.getTableName("task_log") + " (start_time)",
		
		// PostgreSQL特定的优化
		"ANALYZE " + m.getTableName("task"),
		"ANALYZE " + m.getTableName("task_log"),
		"ANALYZE " + m.getTableName("user"),
		"ANALYZE " + m.getTableName("role"),
	}
	
	return m.executeSQLs(sqls)
}

// getTableName 获取带前缀的表名
func (m *MigrationManager) getTableName(tableName string) string {
	prefix := viper.GetString("db.prefix")
	return prefix + tableName
}

// executeSQLs 执行SQL语句列表
func (m *MigrationManager) executeSQLs(sqls []string) error {
	for _, sql := range sqls {
		if err := m.db.Exec(sql).Error; err != nil {
			zap.S().Warnf("执行SQL失败: %s, 错误: %v", sql, err)
			// 继续执行其他SQL，不中断流程
		} else {
			zap.S().Debugf("执行SQL成功: %s", sql)
		}
	}
	return nil
}

// CheckDatabaseVersion 检查数据库版本
func (m *MigrationManager) CheckDatabaseVersion() error {
	var version string
	var sql string
	
	switch m.engine {
	case "mysql":
		sql = "SELECT VERSION()"
	case "postgresql", "postgres":
		sql = "SELECT version()"
	default:
		return nil
	}
	
	err := m.db.Raw(sql).Scan(&version).Error
	if err != nil {
		return err
	}
	
	zap.S().Infof("数据库版本: %s", version)
	return nil
}

// CreateInitialData 创建初始数据
func (m *MigrationManager) CreateInitialData() error {
	zap.S().Info("检查并创建初始数据...")
	
	// 检查是否已有管理员用户
	var count int64
	err := m.db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return err
	}
	
	if count == 0 {
		zap.S().Info("创建默认管理员用户...")
		// 这里可以创建默认的管理员用户
		// 实际实现时需要根据具体需求来创建
	}
	
	return nil
}
