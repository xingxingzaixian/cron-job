package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	taskManager "cronJob/internal/service/cron/task_manager"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
)

type TaskApi struct{}

func TaskRegister(group *gin.RouterGroup) {
	service := &TaskApi{}
	group.GET("/list", service.GetTaskList)
	group.GET("/view", service.GetTaskDetail)
	group.POST("/edit", service.EditTask)
	group.POST("/op", service.OpTask)
}

// GetTaskList godoc
// @Summary 任务列表
// @Description 任务列表
// @Tags 任务管理
// @Security ApiKeyAuth
// @ID /api/task/list
// @Accept json
// @Produce json
// @Param pageNo query int true "页数"
// @Param pageSize query int true "每页条数"
// @Param name query string false "任务名称"
// @Param protocol query string false "任务协议"
// @Param tag query string false "任务标签"
// @Success 200 {object} schemas.Response{data=schemas.SearchTaskResponse} "success"
// @Router /api/task/list [get]
func (s *TaskApi) GetTaskList(ctx *gin.Context) {
	params := &schemas.SearchTaskParmas{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.TaskSearchParamInvalid, err)
		return
	}

	task := &models.Task{}
	tasks, count, err := task.PageList(global.GormDB, params)
	if err != nil {
		schemas.ResponseError(ctx, schemas.TaskSearchListInvalid, err)
		return
	}

	var result []schemas.TaskItemOutput
	for _, v := range tasks {
		result = append(result, schemas.TaskItemOutput{
			TaskEditHTTPInput: schemas.TaskEditHTTPInput{
				ID:            v.ID,
				Name:          v.Name,
				Spec:          v.Spec,
				Protocol:      v.Protocol,
				Command:       v.Command,
				Params:        v.Params,
				Timeout:       v.Timeout,
				Policy:        v.Policy,
				Count:         v.Count,
				Delay:         v.Delay,
				RetryTimes:    v.RetryTimes,
				RetryInterval: v.RetryInterval,
				Tag:           v.Tag,
				Remark:        v.Remark,
				Status:        v.Status,
			},
		})
	}

	schemas.ResponseSuccess(ctx, schemas.SearchTaskResponse{
		Total: count,
		List:  result,
	})
}

// GetTaskDetail godoc
// @Summary 任务信息
// @Description 任务信息
// @Tags 任务管理
// @Security ApiKeyAuth
// @ID /api/task/view
// @Accept json
// @Produce json
// @Param id query int true "任务ID"
// @Success 200 {object} schemas.Response{data=schemas.TaskItemOutput} "success"
// @Router /api/task/view [get]
func (s *TaskApi) GetTaskDetail(ctx *gin.Context) {
	params := &schemas.TaskViewInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.TaskSearchParamInvalid, err)
		return
	}

	task := &models.Task{}
	if err := task.FindOne(global.GormDB, g.Map{
		"id": params.ID,
	}); err != nil {
		schemas.ResponseError(ctx, schemas.TaskNotExist, err)
		return
	}

	schemas.ResponseSuccess(ctx, schemas.TaskItemOutput{
		TaskEditHTTPInput: schemas.TaskEditHTTPInput{
			ID:            task.ID,
			Name:          task.Name,
			Spec:          task.Spec,
			Protocol:      task.Protocol,
			Command:       task.Command,
			Params:        task.Params,
			Timeout:       task.Timeout,
			Policy:        task.Policy,
			Count:         task.Count,
			Delay:         task.Delay,
			RetryTimes:    task.RetryTimes,
			RetryInterval: task.RetryInterval,
			Tag:           task.Tag,
			Remark:        task.Remark,
			Status:        task.Status,
		},
	})
}

// EditTask godoc
// @Summary 新增/修改任务
// @Description 新增/修改任务
// @Tags 任务管理
// @Security ApiKeyAuth
// @ID /api/task/edit
// @Accept json
// @Produce json
// @Param data body schemas.TaskEditHTTPInput true "body"
// @Success 200 {object} schemas.Response{data=string} "success"
// @Router /api/task/edit [post]
func (s *TaskApi) EditTask(ctx *gin.Context) {
	params := &schemas.TaskEditHTTPInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.TaskAddParamInvalid, err)
		return
	}

	task := &models.Task{}
	tx := global.GormDB.Begin()

	if params.ID > 0 {
		err := task.FindOne(tx, g.Map{
			"id": params.ID,
		})
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(ctx, schemas.TaskNotExist, err)
			return
		}
		_, err = task.Update(params.ID, g.Map{
			"name":           params.Name,
			"command":        params.Command,
			"params":         params.Params,
			"spec":           params.Spec,
			"protocol":       params.Protocol,
			"timeout":        params.Timeout,
			"policy":         params.Policy,
			"count":          params.Count,
			"delay":          params.Delay,
			"retry_times":    params.RetryTimes,
			"retry_interval": params.RetryInterval,
			"tag":            params.Tag,
			"remark":         params.Remark,
		})
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(ctx, schemas.TaskUpdateError, err)
			return
		}

		// 如果更新的任务状态是禁用，就删除当前正在调度的任务
		if task.Status == global.TaskStatusDisabled {
			taskManager.TaskManager.RemoveTask(task)
		} else {
			// 添加任务到调度进程中
			taskManager.TaskManager.UpdateTask(task)
		}
	} else {
		if ok := task.IsNameExist(params.Name); ok {
			schemas.ResponseError(ctx, schemas.TaskAlreadyExist, errors.New("任务名已存在"))
			return
		}
		task.Name = params.Name
		task.Command = params.Command
		task.Params = params.Params
		task.Spec = params.Spec
		task.Protocol = params.Protocol
		task.Timeout = params.Timeout
		task.Policy = params.Policy
		task.Count = params.Count
		task.Delay = params.Delay
		task.RetryTimes = params.RetryTimes
		task.RetryInterval = params.RetryInterval
		task.Tag = params.Tag
		task.Remark = params.Remark
		_, err := task.Create()
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(ctx, schemas.TaskCreateError, err)
			return
		}
	}

	tx.Commit()
	schemas.ResponseSuccess(ctx, "")
}

// OpTask godoc
// @Summary 操作任务
// @Description 操作任务
// @Tags 任务管理
// @Security ApiKeyAuth
// @ID /api/task/op
// @Accept json
// @Produce json
// @Param data body schemas.TaskOptionInput true "操作信息"
// @Success 200 {object} schemas.Response{data=bool} "success"
// @Router /api/task/op [post]
func (s *TaskApi) OpTask(ctx *gin.Context) {
	params := schemas.TaskOptionInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.TaskParamInvalid, err)
		return
	}

	tx := global.GormDB.Begin()
	// 查询任务信息
	task := &models.Task{}
	err := task.FindOne(tx, g.Map{
		"id": params.ID,
	})
	if err != nil {
		tx.Rollback()
		schemas.ResponseError(ctx, schemas.TaskNotExist, err)
		return
	}

	switch params.Op {
	case "start":
		_, err = task.Update(params.ID, g.Map{"status": global.TaskStatusEnabled})
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(ctx, schemas.TaskUpdateError, err)
			return
		}
		taskManager.TaskManager.StartTask(task)
	case "stop":
		_, err = task.Update(params.ID, g.Map{"status": global.TaskStatusDisabled})
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(ctx, schemas.TaskUpdateError, err)
			return
		}
		taskManager.TaskManager.StopTask(task)
	case "run":
		taskManager.TaskManager.RunTask(task)
	case "delete":
		task.Delete(tx, params.ID)
		taskManager.TaskManager.RemoveTask(task)
	}

	tx.Commit()
	schemas.ResponseSuccess(ctx, true)
}
