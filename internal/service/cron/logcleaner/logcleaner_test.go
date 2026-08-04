package logcleaner

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	"cronJob/lib/database"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	viper.Set("db.engine", "sqlite")
	viper.Set("db.path", filepath.Join(t.TempDir(), "test_cleaner.db"))

	sqliteDB := &database.SQLiteDB{}
	db, err := sqliteDB.Create(&gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&models.TaskLog{}); err != nil {
		t.Fatalf("迁移任务日志表失败: %v", err)
	}
	return db
}

func insertLog(t *testing.T, taskID uint, start time.Time) {
	t.Helper()
	log := &models.TaskLog{
		TaskId:    taskID,
		TaskName:  "test-task",
		Protocol:  1,
		Status:    global.TaskStatusFinish,
		StartTime: &start,
		EndTime:   &start,
	}
	if _, err := log.Create(); err != nil {
		t.Fatalf("插入日志失败: %v", err)
	}
}

func TestCleanupExpiredLogs(t *testing.T) {
	oldDB := global.GormDB
	defer func() { global.GormDB = oldDB }()

	db := setupTestDB(t)
	global.GormDB = db

	now := time.Now()
	oldTime := now.AddDate(0, 0, -60) // 60天前
	recent := now.AddDate(0, 0, -5)   // 5天前
	insertLog(t, 1, oldTime)
	insertLog(t, 1, oldTime)
	insertLog(t, 2, recent)

	viper.Set("log.retention_days", 30)
	defer viper.Set("log.retention_days", 30)

	deleted, err := CleanupExpiredLogs()
	if err != nil {
		t.Fatalf("清理过期日志失败: %v", err)
	}
	if deleted != 2 {
		t.Fatalf("期望删除2条过期日志, 实际 %d", deleted)
	}

	var remain int64
	if err := db.Model(&models.TaskLog{}).Count(&remain).Error; err != nil {
		t.Fatalf("统计剩余日志失败: %v", err)
	}
	if remain != 1 {
		t.Fatalf("期望剩余1条日志, 实际 %d", remain)
	}
}

func TestCleanupExpiredLogsDisabled(t *testing.T) {
	oldDB := global.GormDB
	defer func() { global.GormDB = oldDB }()

	db := setupTestDB(t)
	global.GormDB = db
	insertLog(t, 1, time.Now().AddDate(0, 0, -60))

	viper.Set("log.retention_days", 0)
	defer viper.Set("log.retention_days", 30)

	deleted, err := CleanupExpiredLogs()
	if err != nil {
		t.Fatalf("清理过期日志失败: %v", err)
	}
	if deleted != 0 {
		t.Fatalf("retention_days=0 时不应删除任何日志, 实际删除 %d", deleted)
	}

	var remain int64
	if err := db.Model(&models.TaskLog{}).Count(&remain).Error; err != nil {
		t.Fatalf("统计剩余日志失败: %v", err)
	}
	if remain != 1 {
		t.Fatalf("期望日志全部保留, 实际剩余 %d", remain)
	}
}

func TestTaskLogPageListDateFilter(t *testing.T) {
	oldDB := global.GormDB
	defer func() { global.GormDB = oldDB }()

	db := setupTestDB(t)
	global.GormDB = db

	now := time.Now()
	oldTime := now.AddDate(0, 0, -30)
	recent := now.AddDate(0, 0, -1)
	insertLog(t, 1, oldTime)
	insertLog(t, 2, recent)

	params := &schemas.TaskLogListInput{
		FormPage:  schemas.FormPage{PageNo: 1, PageSize: 15},
		StartTime: now.AddDate(0, 0, -7).Format("2006-01-02 15:04:05"),
		EndTime:   now.AddDate(0, 0, 1).Format("2006-01-02 15:04:05"),
	}

	logs, _, err := (&models.TaskLog{}).PageList(db, params)
	if err != nil {
		t.Fatalf("分页查询失败: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("时间范围过滤后应返回1条日志, 实际 %d", len(logs))
	}
	if logs[0].TaskId != 2 {
		t.Fatalf("应返回最近的任务日志, 实际 taskId=%d", logs[0].TaskId)
	}
}
