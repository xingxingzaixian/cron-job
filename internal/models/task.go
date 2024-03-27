package models

import (
	"cronJob/internal/global"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Name          string                `gorm:"size:32;not null" json:"name"`
	Level         global.TaskLevel      `gorm:"type:tinyint;not null;default:1" json:"level"`
	Spec          string                `gorm:"size:64;not null" json:"spec"`
	Protocol      global.TaskProtocol   `gorm:"type:tinyint;not null;index" json:"protocol"`
	Command       string                `gorm:"size:256;not null" json:"command"`
	HttpMethod    global.TaskHTTPMethod `gorm:"type:tinyint;not null;default:1" json:"http_method"`
	Timeout       int                   `gorm:"type:mediumint;not null;default:0" json:"timeout"`
	Multi         int8                  `gorm:"type:tinyint;not null;default:1" json:"multi"`
	RetryTimes    int8                  `gorm:"type:tinyint;not null;default:0" json:"retry_times"`
	RetryInterval int16                 `gorm:"type:smallint;not null;default:0" json:"retry_interval"`
	Tag           string                `gorm:"size:32;not null;default:''" json:"tag"`
	Remark        string                `gorm:"size:100;not null;default:''" json:"remark"`
	Status        global.Status         `gorm:"type:tinyint;not null;index;default:0" json:"status"`
	NextRunTime   time.Time             `json:"next_run_time" gorm:"-"`
}

// 获取所有
func (t *Task) GetActiveTasks() (tasks []*Task) {
	global.GormDB.Find(&tasks)
	return
}
