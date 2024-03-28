package handler

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
)

type Handler interface {
	Run(taskModel *models.Task, taskUniqueId uint) (string, error)
}

// HTTPHandler HTTP 任务
type HTTPHandler struct{}

func (h *HTTPHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {
	if taskModel.Timeout <= 0 || taskModel.Timeout > global.HttpExecTimeout {
		taskModel.Timeout = global.HttpExecTimeout
	}
	return "", nil
}

// RPCHandler RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {
	return "", nil
}

// ShellHandler 脚本任务
type SHELLHandler struct{}

func (h *SHELLHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {

	return "", nil
}
