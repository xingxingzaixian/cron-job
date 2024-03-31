package handler

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/utils"
	"fmt"
	"go.uber.org/zap"
	"time"
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
	fmt.Println(22222222222, taskModel.Command)
	return "", nil
}

// ShellHandler 脚本任务
type SHELLHandler struct{}

func (h *SHELLHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(err)
		}
	}()

	zap.S().Infof("execute cmd start: [id: %d cmd: %s]", taskUniqueId, taskModel.Command)
	timeout := global.ShellExecTimeout
	if taskModel.Timeout > 0 {
		timeout = taskModel.Timeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	output, err := utils.ExecShell(ctx, taskModel.Command)
	if err != nil {
		return "", err
	}
	zap.S().Infof("execute cmd start: [id: %d cmd: %s]", taskUniqueId, taskModel.Command)
	return output, nil
}
