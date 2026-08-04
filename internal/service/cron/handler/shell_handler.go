package handler

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/utils"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ShellHandler 脚本任务
type SHELLHandler struct{}

func (h *SHELLHandler) Run(taskModel *models.Task, taskUniqueId uint) (output string, err error) {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("Shell任务执行panic: taskId=%d, panic=%v", taskModel.ID, r)
			err = fmt.Errorf("任务【%d】执行发生内部错误: %v", taskModel.ID, r)
			output = ""
		}
	}()

	zap.S().Infof("execute cmd start: [id: %d cmd: %s]", taskUniqueId, taskModel.Command)
	if taskModel.Timeout <= 0 {
		taskModel.Timeout = global.ShellExecTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taskModel.Timeout)*time.Second)
	defer cancel()

	output, err = utils.ExecShell(ctx, taskModel.Command)
	if err != nil {
		return "", err
	}
	zap.S().Infof("execute cmd finish: [id: %d cmd: %s]", taskUniqueId, taskModel.Command)
	return output, nil
}
