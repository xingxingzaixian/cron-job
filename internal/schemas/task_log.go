package schemas

import (
	"cronJob/internal/global"
	"github.com/gin-gonic/gin"
)

type TaskLogListInput struct {
	FormPage
	TaskID   int               `json:"taskId" form:"taskId" comment:"任务id" example:"1"`
	TaskName string            `json:"taskName" form:"taskName" comment:"任务名称" example:"" `
	Status   global.TaskStatus `json:"status" form:"status" comment:"任务状态" example:"" `
}

func (param *TaskLogListInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type TaskLogInput struct {
	ID int `json:"id" form:"id" comment:"任务日志id" example:"1"`
}

func (param *TaskLogInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type TaskLogItemOutput struct {
	ID         uint                `json:"id" form:"id" comment:"任务日志id"`
	TaskId     uint                `json:"taskId" form:"taskId" comment:"任务id"`         // 任务id
	TaskName   string              `json:"taskName" form:"taskName" comment:"任务名称"`     // 任务名称
	Protocol   global.TaskProtocol `json:"protocol" form:"protocol" comment:"任务协议"`     // 任务协议
	RetryTimes int8                `json:"retryTimes" form:"retryTimes" comment:"重试次数"` // 任务重试次数
	Status     global.TaskStatus   `json:"status" form:"status" comment:"任务状态"`         // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result     string              `json:"result" form:"result" comment:"执行结果"`         // 执行结果
	StartTime  string              `json:"startTime" form:"startTime" comment:"开始时间"`   // 开始时间
	EndTime    string              `json:"endTime" form:"endTime" comment:"结束时间"`       // 结束时间
	TotalTime  int64               `json:"totalTime" form:"totalTime" comment:"执行总时长"`  // 执行总时长
}

type TaskLogOutput struct {
	Total int64               `json:"total" comment:"总数"`
	List  []TaskLogItemOutput `json:"list" comment:"列表"`
}
