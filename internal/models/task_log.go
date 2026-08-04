package models

import (
	"cronJob/internal/global"
	"cronJob/internal/schemas"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

// TaskLog 任务执行日志
type TaskLog struct {
	gorm.Model
	TaskId     uint                `json:"task_id" gorm:"type:int;not null;default:0;index:idx_task_log_task_id"`   // 任务id
	TaskName   string              `json:"task_name" gorm:"type:varchar(32);not null;index:idx_task_log_task_name"` // 任务名称
	Protocol   global.TaskProtocol `json:"protocol" gorm:"type:tinyint;not null;default:1"`                         // 任务方式 1:HTTP  2:间Shell
	RetryTimes int8                `json:"retry_times" gorm:"type:tinyint;not null;default:0"`                      // 任务重试次数
	Status     global.TaskStatus   `json:"status" gorm:"type:tinyint;not null;index:idx_task_log_status;default:0"` // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result     string              `json:"result" gorm:"type:mediumtext"`                                           // 执行结果
	StartTime  *time.Time          `json:"start_time" gorm:"type:datetime;index:idx_task_log_start_time"`           // 开始时间
	EndTime    *time.Time          `json:"end_time" gorm:"type:datetime"`                                           // 结束时间
	Duration   int64               `json:"duration" gorm:"type:bigint;default:0"`                                   // 执行时长(毫秒)
}

// Create 创建任务日志
func (t *TaskLog) Create() (uint, error) {
	result := global.GormDB.Create(t)
	if result.Error != nil {
		return 0, result.Error
	}
	return t.ID, nil
}

// Update 根据ID更新任务日志
func (t *TaskLog) Update(id uint, data g.Map) (int64, error) {
	result := global.GormDB.Model(t).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (t *TaskLog) PageList(tx *gorm.DB, params *schemas.TaskLogListInput) (taskLogs []TaskLog, count int64, err error) {
	// 分页参数归一化：默认值 + 上限
	if params.PageNo <= 0 {
		params.PageNo = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 15
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	query := tx.Model(t)
	if params.TaskID > 0 {
		query = query.Where("task_id = ?", params.TaskID)
	}

	if params.TaskName != "" {
		query = query.Where("task_name like ?", "%"+params.TaskName+"%")
	}

	// 仅当显式传入有效的日志状态（>0）时才过滤；
	// 0/-1 均视为"不限状态"，避免未传状态时误过滤出空结果
	if params.Status > 0 {
		query = query.Where("status = ?", params.Status)
	}

	// 开始时间范围过滤
	if params.StartTime != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", params.StartTime, time.Local); err == nil {
			query = query.Where("start_time >= ?", t)
		}
	}
	if params.EndTime != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", params.EndTime, time.Local); err == nil {
			query = query.Where("start_time <= ?", t)
		}
	}

	query.Count(&count)

	offset := (params.PageNo - 1) * params.PageSize
	result := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&taskLogs)
	if result.Error != nil {
		return
	}

	return
}

func (t *TaskLog) FindOne(tx *gorm.DB, taskLogModel g.Map) error {
	result := tx.Where(taskLogModel).First(t)
	if result.RowsAffected == 0 {
		return errors.New("任务日志不存在")
	}
	return nil
}

// Delete 批量删除任务日志（硬删除）
// 日志属于可清理的审计数据，直接物理删除，避免表无限膨胀
func (t *TaskLog) Delete(tx *gorm.DB, ids []uint) error {
	result := tx.Unscoped().Where("id in (?)", ids).Delete(t)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
