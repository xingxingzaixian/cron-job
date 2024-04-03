package schemas

import (
	"cronJob/internal/global"
	"github.com/gin-gonic/gin"
)

type TaskEditHTTPInput struct {
	ID            uint                `json:"id" form:"id" comment:"任务ID" default:"0"`
	Name          string              `json:"name" form:"name" comment:"任务名称" example:"" validate:"required"`
	Spec          string              `json:"spec" form:"spec" comment:"任务运行时间" example:"" validate:"required"`
	Protocol      global.TaskProtocol `json:"protocol" form:"protocol" comment:"任务协议" example:"1" validate:"required"`
	Command       string              `json:"command" form:"command" comment:"任务命令" example:"" validate:"required"`
	Params        string              `json:"params" form:"params" comment:"命令参数" example:"" `
	Timeout       int                 `json:"timeout" form:"timeout" comment:"超时时间" example:"0" default:"0"`
	Policy        global.TaskPolicy   `json:"policy" form:"policy" comment:"任务策略" example:"1" default:"1"`
	Count         int                 `json:"count" form:"count" comment:"执行次数" example:"0" default:"0"`
	Delay         int                 `json:"delay" form:"delay" comment:"延迟执行" example:"0" default:"0"`
	RetryTimes    int8                `json:"retry_times" form:"retry_times" comment:"重试次数" example:"0" default:"0"`
	RetryInterval int16               `json:"retry_interval" form:"retry_interval" comment:"重试间隔" example:"0" default:"0"`
	Tag           string              `json:"tag" form:"tag" comment:"任务标签" example:"" default:""`
	Remark        string              `json:"remark" form:"remark" comment:"任务备注" example:"" default:""`
	Status        global.TaskStatus   `json:"status" form:"status" comment:"任务状态" default:"0"`
}

func (param *TaskEditHTTPInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type SearchTaskParmas struct {
	FormPage
	Name     string              `json:"name" form:"name" comment:"任务名称" example:""`
	Protocol global.TaskProtocol `json:"protocol" form:"protocol" comment:"任务协议" example:""`
	Tag      string              `json:"tag" form:"tag" comment:"任务标签" example:""`
}

func (param *SearchTaskParmas) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type TaskItemOutput struct {
	TaskEditHTTPInput
}

type SearchTaskResponse struct {
	Total int64            `json:"total" comment:"总数"`
	List  []TaskItemOutput `json:"list" comment:"列表"`
}

type TaskOptionInput struct {
	ID uint   `json:"id" form:"id" comment:"任务ID" example:"1" validate:"required"`
	Op string `json:"op" form:"op" comment:"操作" example:"stop" comment:"start|stop|run|delete" validate:"required"`
}

func (param *TaskOptionInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type TaskViewInput struct {
	ID uint `json:"id" form:"id" comment:"任务ID" example:"1" validate:"required"`
}

func (param *TaskViewInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}
