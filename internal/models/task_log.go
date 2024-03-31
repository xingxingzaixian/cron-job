package models

import (
	"cronJob/internal/global"
	"cronJob/internal/schemas"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

// TaskLog 任务执行日志
type TaskLog struct {
	gorm.Model
	TaskId     uint                `json:"task_id" gorm:"type:int;not null;default 0"`       // 任务id
	Name       string              `json:"name" gorm:"size:32;not null"`                     // 任务名称
	Spec       string              `json:"spec" gorm:"size:64;not null"`                     // crontab
	Protocol   global.TaskProtocol `json:"protocol" gorm:"type:tinyint;not null;index"`      // 协议 1:http 2:RPC
	Command    string              `json:"command" gorm:"size:512;not null"`                 // URL地址或shell命令
	Timeout    int                 `json:"timeout" gorm:"type:mediumint;not null;default:0"` // 任务执行超时时间(单位秒),0不限制
	Policy     global.TaskPolicy   `gorm:"type:tinyint;not null;default:0" json:"policy"`
	RetryTimes int8                `json:"retry_times" gorm:"type:tinyint;not null;default:0"`  // 任务重试次数
	Status     global.TaskStatus   `json:"status" gorm:"type:tinyint;not null;index;default:0"` // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result     string              `json:"result" gorm:"type:mediumtext"`                       // 执行结果
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
	query := tx.Model(t)
	if params.ID > 0 {
		query = query.Where("task_id = ?", params.ID)
	}

	offset := (params.PageNo - 1) * params.PageSize
	result := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&taskLogs)
	if result.Error != nil {
		return
	}
	query.Count(&count)
	return
}

func (t *TaskLog) FindOne(tx *gorm.DB, taskLogModel g.Map) error {
	result := tx.Where(taskLogModel).First(t)
	if result.RowsAffected == 0 {
		return errors.New("任务日志不存在")
	}
	return nil
}
