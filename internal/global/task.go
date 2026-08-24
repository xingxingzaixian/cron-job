package global

type TaskProtocol int8
type TaskStatus int8
type TaskPolicy int8
type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  TaskHTTPMethod = 1
	TaskHttpMethodPost TaskHTTPMethod = 2
)

const (
	TaskProtocolHttp  TaskProtocol = 1
	TaskProtocolShell TaskProtocol = 2
	TaskProtocolSSH   TaskProtocol = 3
)

const (
	TaskPolicyMulti  TaskPolicy = 1 // 并行策略
	TaskPolicyOnce   TaskPolicy = 2 // 单词策略
	TaskPolicySingle TaskPolicy = 3 // 单利策略
	TaskPolicyTimes  TaskPolicy = 4 // 多次策略
)

const (
	TaskStatusDisabled TaskStatus = 0 + iota // 禁用
	TaskStatusEnabled                        // 启用
	TaskStatusFailure                        // 失败
	TaskStatusRunning                        // 运行中
	TaskStatusFinish                         // 成功
	TaskStatusCancel                         // 取消
	TaskStatusTimeout                        // 超时
)

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// NotificationType 通知类型
type NotificationType string

const (
	NotificationTypeEmail   NotificationType = "email"   // 邮件通知
	NotificationTypeWebhook NotificationType = "webhook" // Webhook通知
)

// NotificationTrigger 通知触发条件
type NotificationTrigger string

const (
	NotificationTriggerSuccess NotificationTrigger = "success" // 成功时触发
	NotificationTriggerFailure NotificationTrigger = "failure" // 失败时触发
	NotificationTriggerAll     NotificationTrigger = "all"     // 所有情况触发
	NotificationTriggerTimeout NotificationTrigger = "timeout" // 超时时触发
	NotificationTriggerCancel  NotificationTrigger = "cancel"  // 取消时触发
)

// NotificationStatus 通知状态
type NotificationStatus int8

const (
	NotificationStatusDisabled NotificationStatus = 0 // 禁用
	NotificationStatusEnabled  NotificationStatus = 1 // 启用
)

// NotificationTemplateType 模板类型
type NotificationTemplateType string

const (
	NotificationTemplateTypeSuccess NotificationTemplateType = "success" // 成功模板
	NotificationTemplateTypeFailure NotificationTemplateType = "failure" // 失败模板
)
