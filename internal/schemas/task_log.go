package schemas

import (
	"cronJob/internal/global"
	"github.com/gin-gonic/gin"
)

type TaskLogListInput struct {
	FormPage
	ID int `json:"id" form:"id" comment:"任务id" example:"1"`
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
	TaskId     uint                `json:"task_id" form:"task_id" comment:"任务id"`         // 任务id
	Name       string              `json:"name" form:"name" comment:"任务名称"`               // 任务名称
	Spec       string              `json:"spec" form:"spec" comment:"任务crontab"`          // crontab
	Protocol   global.TaskProtocol `json:"protocol" form:"protocol" comment:"任务协议"`       // 协议 1:http 2:Shell 3:RPC
	Command    string              `json:"command" form:"command" comment:"任务命令"`         // URL地址或shell命令
	Timeout    int                 `json:"timeout" form:"timeout" comment:"超时时间"`         // 任务执行超时时间(单位秒),0不限制
	Policy     global.TaskPolicy   `json:"policy" form:"policy" comment:"任务策略"`           // 任务策略
	RetryTimes int8                `json:"retry_times" form:"retry_times" comment:"重试次数"` // 任务重试次数
	Status     global.TaskStatus   `json:"status" form:"status" comment:"任务状态"`           // 状态 0:执行失败 1:执行中  2:执行完毕 3:任务取消(上次任务未执行完成) 4:异步执行
	Result     string              `json:"result" form:"result" comment:"执行结果"`           // 执行结果
	TotalTime  int64               `json:"total_time" form:"total_time" comment:"执行总时长"`  // 执行总时长
}
type TaskLogOutput struct {
	Total int64               `json:"total" comment:"总数"`
	List  []TaskLogItemOutput `json:"list" comment:"列表"`
}
