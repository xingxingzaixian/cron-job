package job

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/handler"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"go.uber.org/zap"
	"time"
)

func createHandler(taskModel *models.Task) handler.Handler {
	var h handler.Handler = nil
	switch taskModel.Protocol {
	case global.TaskProtocolHttp:
		h = new(handler.HTTPHandler)
	case global.TaskProtocolGrpc:
		h = new(handler.RPCHandler)
	case global.TaskProtocolShell:
		h = new(handler.SHELLHandler)
	}

	return h
}

// CreateJob generates a job function based on the provided task model.
//
// It takes a models.Task parameter and returns a gcron.JobFunc.
// 返回闭包函数，必须传递对象而非指针，否则闭包函数只会在最后一次生效
func CreateJob(taskModel models.Task) gcron.JobFunc {
	hdler := createHandler(&taskModel)
	if hdler == nil {
		return nil
	}

	return func(ctx context.Context) {
		taskLogId := beforeExecJob(&taskModel)
		if taskLogId <= 0 {
			return
		}

		zap.S().Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
		taskResult := execJob(hdler, &taskModel, taskLogId)
		zap.S().Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
		afterExecJob(&taskModel, taskResult, taskLogId)
	}
}

// 任务前置操作
func beforeExecJob(taskModel *models.Task) (taskLogId uint) {
	taskLogId, err := createTaskLog(taskModel, global.TaskStatusRunning)
	if err != nil {
		zap.S().Error("任务开始执行#写入任务日志失败-", err)
		return
	}
	zap.S().Debugf("任务命令-%s", taskModel.Command)

	return taskLogId
}

// 执行任务
func execJob(handler handler.Handler, taskModel *models.Task, taskUniqueId uint) global.TaskResult {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error("panic#service/cron/job/job.go:execJob#", err)
		}
	}()

	// 默认只运行任务一次
	var execTimes int8 = 1
	if taskModel.RetryTimes > 0 {
		execTimes += taskModel.RetryTimes
	}

	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(taskModel, taskUniqueId)
		if err == nil {
			return global.TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++

		if i < execTimes {
			zap.S().Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", taskModel.ID, i, output, err.Error())
			if taskModel.RetryInterval > 0 {
				time.Sleep(time.Duration(taskModel.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}

	return global.TaskResult{Result: output, Err: err, RetryTimes: taskModel.RetryTimes}
}

// 任务执行后置操作
func afterExecJob(taskModel *models.Task, taskResult global.TaskResult, taskLogId uint) {
	_, err := updateTaskLog(taskLogId, taskResult)
	if err != nil {
		zap.S().Error("任务结束#更新任务日志失败-", err)
	}
}

func createTaskLog(taskModel *models.Task, status global.TaskStatus) (insertId uint, err error) {
	fmt.Printf("1111111111111111111111: %+v\n", taskModel)
	taskLogModel := new(models.TaskLog)
	taskLogModel.TaskId = taskModel.ID
	taskLogModel.TaskName = taskModel.Name
	taskLogModel.Protocol = taskModel.Protocol
	taskLogModel.RetryTimes = taskModel.RetryTimes
	taskLogModel.Status = status
	insertId, err = taskLogModel.Create()
	return
}

// 更新任务日志
func updateTaskLog(taskLogId uint, taskResult global.TaskResult) (int64, error) {
	taskLogModel := new(models.TaskLog)
	var status global.TaskStatus
	var result string
	if taskResult.Err != nil {
		status = global.TaskStatusFailure
		result = taskResult.Err.Error()
	} else {
		status = global.TaskStatusFinish
		result = taskResult.Result
	}

	return taskLogModel.Update(taskLogId, g.Map{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
	})
}
