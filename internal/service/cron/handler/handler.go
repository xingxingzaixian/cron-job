package handler

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/lib/httpclient"
	"cronJob/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go.uber.org/zap"
	"time"
)

type Handler interface {
	Run(taskModel *models.Task, taskUniqueId uint) (string, error)
}

// HTTPHandler HTTP 任务
type HTTPHandler struct{}

func (h *HTTPHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(err)
		}
	}()

	if taskModel.Timeout <= 0 {
		taskModel.Timeout = global.HttpExecTimeout
	}

	command := g.Map{}
	err := json.Unmarshal([]byte(taskModel.Command), &command)
	if err != nil {
		return "", fmt.Errorf("任务【%d】参数解析失败-%s", taskModel.ID, err.Error())
	}

	method, ok := command["method"]
	if !ok {
		return "", fmt.Errorf("任务【%d】参数解析失败", taskModel.ID)
	}

	url, ok := command["url"]
	if !ok {
		return "", fmt.Errorf("任务【%d】参数解析失败", taskModel.ID)
	}

	var output string
	zap.S().Infof("execute http start: [id: %d url: %s method: %s]", taskUniqueId, url, method)
	if method.(string) == "GET" {
		output, err = httpclient.Get(url.(string), taskModel.Params, time.Duration(taskModel.Timeout)*time.Second)
	} else if method.(string) == "POST" {
		output, err = httpclient.Post(url.(string), taskModel.Params, time.Duration(taskModel.Timeout)*time.Second)
	} else {
		err = fmt.Errorf("任务【%d】不支持的请求方式【%s】", taskModel.ID, method)
	}
	zap.S().Infof("execute http finish: [id: %d url: %s method: %s]", taskUniqueId, url, method)
	if err != nil {
		return "", err
	}
	return output, nil
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
