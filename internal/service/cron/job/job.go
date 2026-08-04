package job

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/handler"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"go.uber.org/zap"
)

// maxLogOutputSize 单条任务日志内容大小上限（1MB），防止大输出拖垮数据库
const maxLogOutputSize = 1 * 1024 * 1024

func createHandler(taskModel *models.Task) handler.Handler {
	var h handler.Handler = nil
	switch taskModel.Protocol {
	case global.TaskProtocolHttp:
		h = new(handler.HTTPHandler)
	case global.TaskProtocolShell:
		h = new(handler.SHELLHandler)
	case global.TaskProtocolSSH:
		h = new(handler.SSHHandler)
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
		// 检查任务是否有依赖
		hasDependencies, err := taskModel.HasDependencies()
		if err != nil {
			zap.S().Errorf("检查任务依赖失败: taskId=%d, error=%v", taskModel.ID, err)
			return
		}

		// 如果有依赖，检查依赖是否满足
		if hasDependencies {
			canExecute, err := taskModel.CanExecute()
			if err != nil {
				zap.S().Errorf("检查任务执行条件失败: taskId=%d, error=%v", taskModel.ID, err)
				return
			}

			if !canExecute {
				zap.S().Infof("任务依赖未满足，跳过执行: taskId=%d", taskModel.ID)
				// 记录日志，说明依赖未满足；仅当最近一次日志不是取消状态时写入，避免日志膨胀
				lastLog := &models.TaskLog{}
				lastLogErr := global.GormDB.Where("task_id = ?", taskModel.ID).Order("id DESC").First(lastLog).Error
				if lastLogErr != nil || lastLog.Status != global.TaskStatusCancel {
					_, err := createTaskLog(&taskModel, global.TaskStatusCancel)
					if err != nil {
						zap.S().Errorf("创建依赖未满足任务日志失败: taskId=%d, error=%v", taskModel.ID, err)
					}
				}
				return
			}
		}

		startTime := time.Now()
		taskLogId := beforeExecJob(&taskModel)
		if taskLogId <= 0 {
			return
		}

		zap.S().Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
		taskResult := execJob(ctx, hdler, &taskModel, taskLogId)
		zap.S().Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
		afterExecJob(&taskModel, taskResult, taskLogId, startTime)
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
func execJob(ctx context.Context, handler handler.Handler, taskModel *models.Task, taskUniqueId uint) global.TaskResult {
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

		if i >= execTimes {
			break
		}

		interval := retryInterval(taskModel.RetryInterval, i)

		zap.S().Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", taskModel.ID, i, output, err.Error())

		// 可取消的重试等待，任务被取消/系统关闭时立即退出
		select {
		case <-ctx.Done():
			return global.TaskResult{Result: output, Err: fmt.Errorf("任务【%d】执行被取消: %v", taskModel.ID, ctx.Err()), RetryTimes: i}
		case <-time.After(interval):
		}
	}

	return global.TaskResult{Result: output, Err: err, RetryTimes: taskModel.RetryTimes}
}

// retryInterval 计算重试等待间隔（纯函数，便于单测）
// 任务配置的间隔优先；未配置时按重试次数递增，但上限5分钟，避免长时间阻塞
func retryInterval(retryInterval int16, retryTimes int8) time.Duration {
	interval := time.Duration(retryInterval) * time.Second
	if retryInterval <= 0 {
		interval = time.Duration(retryTimes) * time.Minute
		if interval > 5*time.Minute {
			interval = 5 * time.Minute
		}
	}
	return interval
}

// 任务执行后置操作
func afterExecJob(taskModel *models.Task, taskResult global.TaskResult, taskLogId uint, startTime time.Time) {
	_, err := updateTaskLog(taskLogId, taskResult, startTime)
	if err != nil {
		zap.S().Error("任务结束#更新任务日志失败-", err)
	}
}

func createTaskLog(taskModel *models.Task, status global.TaskStatus) (insertId uint, err error) {
	now := time.Now()
	taskLogModel := new(models.TaskLog)
	taskLogModel.TaskId = taskModel.ID
	taskLogModel.TaskName = taskModel.Name
	taskLogModel.Protocol = taskModel.Protocol
	taskLogModel.RetryTimes = taskModel.RetryTimes
	taskLogModel.Status = status
	taskLogModel.StartTime = &now
	insertId, err = taskLogModel.Create()
	if err != nil {
		zap.S().Errorf("创建任务日志失败: taskId=%d, error=%v", taskModel.ID, err)
	}
	return
}

// 更新任务日志
func updateTaskLog(taskLogId uint, taskResult global.TaskResult, startTime time.Time) (int64, error) {
	taskLogModel := new(models.TaskLog)
	var status global.TaskStatus
	var result string
	if taskResult.Err != nil {
		status = global.TaskStatusFailure
		result = truncateLogOutput(taskResult.Err.Error())
	} else {
		status = global.TaskStatusFinish
		result = truncateLogOutput(taskResult.Result)
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime).Milliseconds()

	return taskLogModel.Update(taskLogId, g.Map{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
		"end_time":    endTime,
		"duration":    duration,
	})
}

// truncateLogOutput 截断过长的日志输出，并在末尾标注原始长度
func truncateLogOutput(output string) string {
	if len(output) <= maxLogOutputSize {
		return output
	}
	truncated := output[:maxLogOutputSize]
	return truncated + fmt.Sprintf("\n...[输出已截断，共 %d 字节]", len(output))
}
