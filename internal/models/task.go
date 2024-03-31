package models

import (
	"cronJob/internal/global"
	"cronJob/internal/schemas"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name          string              `gorm:"size:32;not null" json:"name"`
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
	Status        global.TaskStatus   `gorm:"type:tinyint;not null;index;default:1" json:"status"`
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

// 获取所有
func (t *Task) GetActiveTasks() (tasks []Task) {
	global.GormDB.Where("status != ?", global.TaskStatusDisabled).Find(&tasks)
	return
}

// PageList 分页查询
func (t *Task) PageList(tx *gorm.DB, params *schemas.SearchTaskParmas) (tasks []Task, count int64, err error) {
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

	offset := (params.PageNo - 1) * params.PageSize
	result := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&tasks)
	if result.Error != nil {
		return
	}
	query.Count(&count)
	return
}

// IsNameExist 判断任务名是否存在
func (t *Task) IsNameExist(name string) bool {
	var count int64 = 0
	global.GormDB.Model(t).Where("name = ?", name).Count(&count)
	return count > 0
}

func (t *Task) IsExist(id uint) bool {
	var count int64 = 0
	if id > 0 {
		global.GormDB.Model(t).Where("id != ?", id).Count(&count)
		return count > 0
	}
	return false
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

func (t *Task) Delete(tx *gorm.DB, id uint) {
	tx.Delete(&Task{}, id)
}
