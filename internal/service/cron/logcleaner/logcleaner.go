package logcleaner

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/lib/database"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// cleanupBatchSize 单批删除的最大日志条数，分批删除避免大表一次性删除长时间锁表
const cleanupBatchSize = 5000

var (
	lastVacuumMu sync.Mutex
	lastVacuum   time.Time
)

// Start 启动日志自动清理协程：立即执行一次，之后按配置间隔定期执行
func Start() {
	go func() {
		cleanupOnce()
		ticker := time.NewTicker(cleanupInterval())
		defer ticker.Stop()
		for range ticker.C {
			cleanupOnce()
		}
	}()
}

// cleanupInterval 清理检查间隔，默认6小时
func cleanupInterval() time.Duration {
	hours := viper.GetInt("log.cleanup_interval_hours")
	if hours <= 0 {
		hours = 6
	}
	return time.Duration(hours) * time.Hour
}

func cleanupOnce() {
	deleted, err := CleanupExpiredLogs()
	if err != nil {
		zap.S().Errorf("清理过期任务日志失败: %v", err)
	} else if deleted > 0 {
		zap.S().Infof("清理过期任务日志 %d 条", deleted)
	}
	maybeVacuum()
}

// CleanupExpiredLogs 分批删除超过保留期的任务日志
// 返回删除的日志条数；log.retention_days <= 0 表示不自动清理
func CleanupExpiredLogs() (int64, error) {
	retentionDays := viper.GetInt("log.retention_days")
	if retentionDays <= 0 {
		return 0, nil
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	var total int64
	for {
		var ids []uint
		if err := global.GormDB.Model(&models.TaskLog{}).
			Where("start_time < ?", cutoff).
			Order("id ASC").
			Limit(cleanupBatchSize).
			Pluck("id", &ids).Error; err != nil {
			return total, err
		}
		if len(ids) == 0 {
			break
		}

		// 硬删除，真正释放存储空间
		result := global.GormDB.Unscoped().Where("id IN (?)", ids).Delete(&models.TaskLog{})
		if result.Error != nil {
			return total, result.Error
		}
		total += result.RowsAffected
	}
	return total, nil
}

// maybeVacuum 按配置间隔对SQLite执行VACUUM，回收删除后的文件空间
func maybeVacuum() {
	interval := viper.GetInt("db.sqlite.vacuum_interval_hours")
	if interval <= 0 {
		interval = 168
	}

	lastVacuumMu.Lock()
	defer lastVacuumMu.Unlock()
	if time.Since(lastVacuum) < time.Duration(interval)*time.Hour {
		return
	}

	if err := database.VacuumDB(); err != nil {
		zap.S().Errorf("数据库VACUUM失败: %v", err)
		return
	}
	lastVacuum = time.Now()
	zap.S().Info("数据库VACUUM完成")
}
