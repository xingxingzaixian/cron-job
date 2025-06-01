package handler

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/utils"
	"time"

	"go.uber.org/zap"
)

// ShellHandler 脚本任务
type SHELLHandler struct{}

func (h *SHELLHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(err)
		}
	}()

	zap.S().Infof("execute cmd start: [id: %d cmd: %s]", taskUniqueId, taskModel.Command)
	if taskModel.Timeout <= 0 {
		taskModel.Timeout = global.ShellExecTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taskModel.Timeout)*time.Second)
	defer cancel()

	output, err := utils.ExecShell(ctx, taskModel.Command)
	if err != nil {
		return "", err
	}
	zap.S().Infof("execute cmd finish: [id: %d cmd: %s]", taskUniqueId, taskModel.Command)
	return output, nil
}
