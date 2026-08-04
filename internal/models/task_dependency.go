package models

import (
	"cronJob/internal/global"
	"github.com/gogf/gf/v2/frame/g"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TaskDependency 任务依赖关系
type TaskDependency struct {
	gorm.Model
	TaskID      uint              `json:"task_id" gorm:"type:int;not null;index:idx_task_dep_task_id"`           // 任务ID
	DependentID uint              `json:"dependent_id" gorm:"type:int;not null;index:idx_task_dep_dependent_id"` // 依赖的任务ID
	IsMust      bool              `json:"is_must" gorm:"type:tinyint;not null;default:1"`                        // 是否必须依赖（true: 必须成功完成 false: 完成即可）
	Status      global.TaskStatus `json:"status" gorm:"type:tinyint;not null;default:0"`                         // 依赖状态 0:未完成 1:已完成 2:失败
}

// Create 创建任务依赖关系
func (t *TaskDependency) Create() (uint, error) {
	result := global.GormDB.Create(t)
	if result.Error != nil {
		return 0, result.Error
	}
	return t.ID, nil
}

// Update 更新任务依赖关系
func (t *TaskDependency) Update(id uint, data g.Map) (int64, error) {
	result := global.GormDB.Model(t).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// Delete 删除任务依赖关系
func (t *TaskDependency) Delete(id uint) error {
	result := global.GormDB.Delete(t, id)
	return result.Error
}

// DeleteByTaskID 根据任务ID删除依赖关系
func (t *TaskDependency) DeleteByTaskID(taskID uint) error {
	result := global.GormDB.Where("task_id = ?", taskID).Delete(t)
	return result.Error
}

// DeleteByDependentID 根据被依赖任务ID删除依赖关系
func (t *TaskDependency) DeleteByDependentID(dependentID uint) error {
	result := global.GormDB.Where("dependent_id = ?", dependentID).Delete(t)
	return result.Error
}

// GetDependenciesByTaskID 获取任务的所有依赖
func (t *TaskDependency) GetDependenciesByTaskID(taskID uint) ([]TaskDependency, error) {
	var dependencies []TaskDependency
	result := global.GormDB.Where("task_id = ?", taskID).Find(&dependencies)
	if result.Error != nil {
		return nil, result.Error
	}
	return dependencies, nil
}

// GetDependentsByTaskID 获取依赖于指定任务的所有任务
func (t *TaskDependency) GetDependentsByTaskID(taskID uint) ([]TaskDependency, error) {
	var dependencies []TaskDependency
	result := global.GormDB.Where("dependent_id = ?", taskID).Find(&dependencies)
	if result.Error != nil {
		return nil, result.Error
	}
	return dependencies, nil
}

// CheckDependenciesCompleted 检查任务的所有依赖是否已完成
func (t *TaskDependency) CheckDependenciesCompleted(taskID uint) (bool, error) {
	// 获取任务的所有依赖
	dependencies, err := t.GetDependenciesByTaskID(taskID)
	if err != nil {
		return false, err
	}

	// 检查每个依赖是否已完成
	for _, dep := range dependencies {
		// 如果是必须依赖，则需要检查是否成功完成
		if dep.IsMust {
			// 查询依赖任务的最新日志状态
			taskLog := TaskLog{}
			result := global.GormDB.Where("task_id = ?", dep.DependentID).Order("id DESC").First(&taskLog)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					// 没有执行记录，依赖未完成
					zap.S().Infof("Dependency task %d not found for task %d", dep.DependentID, taskID)
					return false, nil
				}
				zap.S().Errorf("Failed to get latest log for task %d: %v", dep.DependentID, result.Error)
				return false, result.Error
			}

			// 检查任务是否成功完成
			if taskLog.Status != global.TaskStatusFinish {
				zap.S().Infof("Dependency task %d not finished for task %d", dep.DependentID, taskID)
				return false, nil
			}
		} else {
			// 非必须依赖，只需要检查是否有执行记录
			var count int64
			global.GormDB.Model(&TaskLog{}).Where("task_id = ?", dep.DependentID).Count(&count)
			if count == 0 {
				zap.S().Infof("No execution record found for non-required dependency task %d", dep.DependentID)
				return false, nil
			}
		}
	}

	return true, nil
}
