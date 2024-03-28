package global

type TaskProtocol int8
type TaskStatus int8
type TaskLevel int8
type TaskPolicy int8
type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  TaskHTTPMethod = 1
	TaskHttpMethodPost TaskHTTPMethod = 2
)

const (
	TaskProtocolHttp  TaskProtocol = 1
	TaskProtocolShell TaskProtocol = 2
	TaskProtocolGrpc  TaskProtocol = 3
)

const (
	TaskPolicyMulti  TaskPolicy = 1 // 并行策略
	TaskPolicyOnce   TaskPolicy = 2 // 单词策略
	TaskPolicySingle TaskPolicy = 3 // 单利策略
	TaskPolicyTimes  TaskPolicy = 4 // 多次策略
)

const (
	TaskStatusFailure TaskStatus = 1 // 失败
	TaskStatusRunning TaskStatus = 2 // 运行中
	TaskStatusFinish  TaskStatus = 3 // 成功
	TaskStatusCancel  TaskStatus = 4 // 取消
	TaskStatusTimeout TaskStatus = 5 // 超时
)

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}
