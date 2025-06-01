package handler

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/lib/httpclient"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"go.uber.org/zap"
)

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
	method = strings.ToUpper(method.(string))
	zap.S().Infof("execute http start: [id: %d url: %s method: %s]", taskUniqueId, url, method)
	if method == "GET" {
		output, err = httpclient.Get(url.(string), taskModel.Params, time.Duration(taskModel.Timeout)*time.Second)
	} else if method == "POST" {
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
