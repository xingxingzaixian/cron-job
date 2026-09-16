package models

import (
	"cronJob/internal/global"
	"cronJob/internal/schemas"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

// 列类型注意：不要使用 tinyint/mediumint/datetime 等 MySQL 专有类型标签，
// GORM 会把 type 标签原样写入 DDL，PostgreSQL 会报 type "xxx" does not exist。
// 可移植写法：int8 自动映射为 MySQL tinyint / PG smallint；
// size:16777215 映射为 MySQL mediumtext / PG text；时间字段不加 type（MySQL datetime(3) / PG timestamptz）。
type Task struct {
	gorm.Model
	Name               string              `gorm:"size:32;not null;index:idx_task_name" json:"name"`
	Spec               string              `gorm:"size:64;not null" json:"spec"`
	Protocol           global.TaskProtocol `gorm:"not null;index:idx_task_protocol" json:"protocol"`
	Command            string              `gorm:"size:512;not null" json:"command"`
	Params             string              `gorm:"size:16777215" json:"params"`
	Timeout            int                 `gorm:"type:integer;not null;default:0" json:"timeout"`
	Policy             global.TaskPolicy   `gorm:"not null;default:1" json:"policy"`
	Count              int                 `gorm:"type:smallint;not null;default:0" json:"count"`
	ExecutedTimes      int                 `gorm:"type:int;not null;default:0" json:"executed_times"` // once/times策略已执行次数（持久化，避免重启后重复执行）
	Delay              int                 `gorm:"type:smallint;not null;default:0" json:"delay"`
	RetryTimes         int8                `gorm:"not null;default:0" json:"retry_times"`
	RetryInterval      int16               `gorm:"type:smallint;not null;default:0" json:"retry_interval"`
	Tag                string              `gorm:"size:32;not null;default:'';index:idx_task_tag" json:"tag"`
	Remark             string              `gorm:"size:256;not null;default:''" json:"remark"`
	Status             global.TaskStatus   `gorm:"not null;default:0;index:idx_task_status" json:"status"`
	EnableNotification bool                `gorm:"type:boolean;not null;default:true" json:"enable_notification"` // 是否启用通知
}

func (t *Task) Create() (uint, error) {
	result := global.GormDB.Create(t)
	if result.Error != nil {
		return 0, result.Error
	}
	return t.ID, nil
}

// Update 根据ID更新任务日志
func (t *Task) Update(id uint, data g.Map) (int64, error) {
	result := global.GormDB.Model(t).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// GetActiveTasks 获取所有
func (t *Task) GetActiveTasks() (tasks []Task) {
	global.GormDB.Where("status != ?", global.TaskStatusDisabled).Find(&tasks)
	return
}

// PageList 分页查询
func (t *Task) PageList(tx *gorm.DB, params *schemas.SearchTaskParmas) (tasks []Task, count int64, err error) {
	// 分页参数归一化：默认值 + 上限，避免拉全表或非法偏移
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
	if params.Name != "" {
		query = query.Where("name like ?", "%"+params.Name+"%")
	}

	if params.Tag != "" {
		query = query.Where("tag = ?", params.Tag)
	}

	if params.Protocol != 0 {
		query = query.Where("protocol = ?", params.Protocol)
	}

	query.Count(&count)

	offset := (params.PageNo - 1) * params.PageSize
	result := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&tasks)
	if result.Error != nil {
		return
	}

	return
}

// IsNameExist 判断任务名是否存在
func (t *Task) IsNameExist(name string) bool {
	var count int64 = 0
	global.GormDB.Model(t).Where("name = ?", name).Count(&count)
	return count > 0
}

func (t *Task) Find(tx *gorm.DB, taskModel g.Map) (list []Task, err error) {
	result := tx.Where(taskModel).Find(&list)
	if result.RowsAffected == 0 {
		return nil, errors.New("任务不存在")
	}

	return
}

func (t *Task) FindOne(tx *gorm.DB, taskModel g.Map) error {
	result := tx.Where(taskModel).First(t)
	if result.RowsAffected == 0 {
		return errors.New("任务不存在")
	}
	return nil
}

func (t *Task) Delete(tx *gorm.DB, id uint) error {
	result := tx.Delete(&Task{}, id)
	return result.Error
}

func (t *Task) GetDependents() ([]TaskDependency, error) {
	dependency := TaskDependency{}
	return dependency.GetDependentsByTaskID(t.ID)
}

// HasDependencies 检查任务是否有依赖
func (t *Task) HasDependencies() (bool, error) {
	dependency := TaskDependency{}
	dependencies, err := dependency.GetDependenciesByTaskID(t.ID)
	if err != nil {
		return false, err
	}
	return len(dependencies) > 0, nil
}

// CanExecute 检查任务是否可以执行（依赖是否满足）
func (t *Task) CanExecute() (bool, error) {
	dependency := TaskDependency{}
	completed, err := dependency.CheckDependenciesCompleted(t.ID)
	if err != nil {
		return false, err
	}
	return completed, nil
}
