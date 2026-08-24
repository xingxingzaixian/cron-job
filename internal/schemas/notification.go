package schemas

import (
	"cronJob/internal/global"
	"github.com/gin-gonic/gin"
)

// NotificationConfigInput 创建/更新通知配置的输入参数
type NotificationConfigInput struct {
	ID       uint                          `json:"id" form:"id" comment:"通知配置ID" example:"0"`
	TaskID   uint                          `json:"task_id" form:"task_id" comment:"任务ID" example:"0"`
	Name     string                        `json:"name" form:"name" comment:"通知名称" example:"任务完成通知" validate:"required"`
	Type     global.NotificationType       `json:"type" form:"type" comment:"通知类型" example:"email" validate:"required"`
	Target   string                        `json:"target" form:"target" comment:"目标地址" example:"user@example.com" validate:"required"`
	Trigger  global.NotificationTrigger    `json:"trigger" form:"trigger" comment:"触发条件" example:"all" default:"all"`
	Enabled  bool                          `json:"enabled" form:"enabled" comment:"是否启用" default:"true"`

	// 邮件配置
	EmailRecipients string `json:"email_recipients" form:"email_recipients" comment:"收件人列表"`
	EmailSubject    string `json:"email_subject" form:"email_subject" comment:"邮件主题模板"`
	EmailBody       string `json:"email_body" form:"email_body" comment:"邮件正文模板"`

	// Webhook配置
	WebhookURL      string `json:"webhook_url" form:"webhook_url" comment:"Webhook URL"`
	WebhookMethod   string `json:"webhook_method" form:"webhook_method" comment:"HTTP方法" default:"POST"`
	WebhookHeaders  string `json:"webhook_headers" form:"webhook_headers" comment:"自定义请求头"`
	WebhookBody     string `json:"webhook_body" form:"webhook_body" comment:"自定义请求体模板"`

	// 重试配置
	RetryTimes    int8  `json:"retry_times" form:"retry_times" comment:"重试次数" default:"0"`
	RetryInterval int   `json:"retry_interval" form:"retry_interval" comment:"重试间隔(秒)" default:"5"`
}

func (param *NotificationConfigInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

// NotificationConfigOutput 通知配置输出
type NotificationConfigOutput struct {
	ID        uint                       `json:"id"`
	TaskID    uint                       `json:"task_id"`
	Name      string                     `json:"name"`
	Type      global.NotificationType    `json:"type"`
	Target    string                     `json:"target"`
	Trigger   global.NotificationTrigger `json:"trigger"`
	Enabled   bool                       `json:"enabled"`
	CreatedAt string                     `json:"created_at"`
	UpdatedAt string                     `json:"updated_at"`
}

// SearchNotificationParams 搜索通知配置参数
type SearchNotificationParams struct {
	FormPage
	TaskID uint `json:"task_id" form:"task_id" comment:"任务ID"`
}

func (param *SearchNotificationParams) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

// NotificationLogInput 查询通知日志参数
type NotificationLogInput struct {
	FormPage
	TaskID         uint `json:"task_id" form:"task_id" comment:"任务ID" validate:"required"`
	NotificationID uint `json:"notification_id" form:"notification_id" comment:"通知配置ID"`
}

func (param *NotificationLogInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

// NotificationLogOutput 通知日志输出
type NotificationLogOutput struct {
	ID             uint   `json:"id"`
	NotificationID uint   `json:"notification_id"`
	TaskID         uint   `json:"task_id"`
	TaskLogID      uint   `json:"task_log_id"`
	Status         int8   `json:"status"`
	Result         string `json:"result"`
	RetryCount     int8   `json:"retry_count"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
}

// NotificationTestInput 测试通知发送参数
type NotificationTestInput struct {
	Type    global.NotificationType `json:"type" form:"type" comment:"通知类型" example:"email" validate:"required"`
	Target  string                  `json:"target" form:"target" comment:"目标地址" example:"user@example.com" validate:"required"`
	Content string                  `json:"content" form:"content" comment:"测试内容" example:"这是一条测试消息"`
}

func (param *NotificationTestInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

// NotificationTestOutput 测试通知发送输出
type NotificationTestOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
