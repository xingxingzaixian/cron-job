package models

import (
	"cronJob/internal/global"
	"github.com/gogf/gf/v2/container/gmap"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Name          string              `gorm:"size:32;not null" json:"name"`
	Level         global.TaskLevel    `gorm:"type:tinyint;not null;default:1" json:"level"`
	Spec          string              `gorm:"size:64;not null" json:"spec"`
	Protocol      global.TaskProtocol `gorm:"type:tinyint;not null;index" json:"protocol"`
	Command       string              `gorm:"size:512;not null" json:"command"`
	Timeout       int                 `gorm:"type:mediumint;not null;default:0" json:"timeout"`
	Policy        global.TaskPolicy   `gorm:"type:tinyint;not null;default:1" json:"policy"`
	Count         int                 `gorm:"type:smallint;not null;default:0" json:"count"`
	Delay         int                 `gorm:"type:smallint;not null;default:0" json:"delay"`
	RetryTimes    int8                `gorm:"type:tinyint;not null;default:0" json:"retry_times"`
	RetryInterval int16               `gorm:"type:smallint;not null;default:0" json:"retry_interval"`
	Tag           string              `gorm:"size:32;not null;default:''" json:"tag"`
	Remark        string              `gorm:"size:100;not null;default:''" json:"remark"`
	Status        global.TaskStatus   `gorm:"type:tinyint;not null;index;default:0" json:"status"`
	NextRunTime   time.Time           `json:"next_run_time" gorm:"-"`
}

func (t *Task) Create() (uint, error) {
	result := global.GormDB.Create(t)
	if result.Error != nil {
		return 0, result.Error
	}
	return t.ID, nil
}

// Update 根据ID更新任务日志
func (t *Task) Update(id uint, data *gmap.Map) (int64, error) {
	result := global.GormDB.Model(t).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// 获取所有
func (t *Task) GetActiveTasks() (tasks []*Task) {
	global.GormDB.Find(&tasks)
	return
}
