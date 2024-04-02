package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
)

type TaskLogApi struct{}

func TaskLogRegister(group *gin.RouterGroup) {
	service := &TaskLogApi{}
	group.GET("/list", service.GetTaskLogList)
	group.GET("/query", service.GetTaskLog)
}

// GetTaskLogList godoc
// @Summary 任务执行日志列表
// @Description 任务执行日志列表
// @Tags 任务日志管理
// @Security ApiKeyAuth
// @ID /api/taskLog/list
// @Accept json
// @Produce json
// @Param pageNo query int true "页数"
// @Param pageSize query int true "每页条数"
// @Param taskId query int true "任务ID"
// @Success 200 {object} schemas.Response{data=schemas.TaskLogOutput} "success"
// @Router /api/taskLog/list [get]
func (s *TaskLogApi) GetTaskLogList(ctx *gin.Context) {
	params := schemas.TaskLogListInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.TaskLogParamInvalid, err)
		return
	}

	taskLog := &models.TaskLog{}
	taskLogs, count, err := taskLog.PageList(global.GormDB, &params)
	if err != nil {
		schemas.ResponseError(ctx, schemas.TaskLogListInvalid, err)
		return
	}

	var result []schemas.TaskLogItemOutput
	for _, v := range taskLogs {
		result = append(result, schemas.TaskLogItemOutput{
			ID:         v.ID,
			TaskId:     v.TaskId,
			RetryTimes: v.RetryTimes,
			Status:     v.Status,
			Result:     v.Result,
			TotalTime:  v.UpdatedAt.Unix() - v.CreatedAt.Unix(),
		})
	}
	schemas.ResponseSuccess(ctx, schemas.TaskLogOutput{
		Total: count,
		List:  result,
	})
}

// GetTaskLog godoc
// @Summary 任务执行日志
// @Description 任务执行日志
// @Tags 任务日志管理
// @Security ApiKeyAuth
// @ID /api/taskLog/query
// @Accept json
// @Produce json
// @Param id query int true "任务日志ID"
// @Success 200 {object} schemas.Response{data=schemas.TaskLogItemOutput} "success"
// @Router /api/taskLog/query [get]
func (s *TaskLogApi) GetTaskLog(ctx *gin.Context) {
	params := schemas.TaskLogInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.TaskLogParamInvalid, err)
		return
	}

	taskLog := &models.TaskLog{}
	err := taskLog.FindOne(global.GormDB, g.Map{"id": params.ID})
	if err != nil {
		schemas.ResponseError(ctx, schemas.TaskLogListInvalid, err)
		return
	}
	schemas.ResponseSuccess(ctx, schemas.TaskLogItemOutput{
		ID:         taskLog.ID,
		TaskId:     taskLog.TaskId,
		RetryTimes: taskLog.RetryTimes,
		Status:     taskLog.Status,
		Result:     taskLog.Result,
		TotalTime:  taskLog.UpdatedAt.Unix() - taskLog.CreatedAt.Unix(),
	})
}
