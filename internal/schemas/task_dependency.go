package schemas

import "cronJob/internal/global"

type TaskDependency struct {
	TaskID       uint              `json:"task_id" form:"task_id" validate:"required"`
	DependentID  uint              `json:"dependent_id" form:"dependent_id" validate:"required"`
	IsMust       bool              `json:"is_must" form:"is_must"`
}

type TaskDependencyList struct {
	TaskID uint `json:"task_id" form:"task_id" validate:"required"`
}

type TaskDependencyResponse struct {
	ID           uint              `json:"id"`
	TaskID       uint              `json:"task_id"`
	DependentID  uint              `json:"dependent_id"`
	IsMust       bool              `json:"is_must"`
	Status       global.TaskStatus `json:"status"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	
	// 关联的任务信息
	DependentTaskName string        `json:"dependent_task_name"`
}