package models

import (
	"cronJob/internal/global"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

// 使用 global 包中定义的类型别名
type NotificationType = global.NotificationType
type NotificationTrigger = global.NotificationTrigger
type NotificationTemplateType = global.NotificationTemplateType

// 重新导出常量以保持向后兼容
const (
	NotificationTypeEmail   = global.NotificationTypeEmail
	NotificationTypeWebhook = global.NotificationTypeWebhook
)

const (
	NotificationTriggerSuccess = global.NotificationTriggerSuccess
	NotificationTriggerFailure = global.NotificationTriggerFailure
	NotificationTriggerAll     = global.NotificationTriggerAll
	NotificationTriggerTimeout = global.NotificationTriggerTimeout
	NotificationTriggerCancel  = global.NotificationTriggerCancel
)

const (
	NotificationTemplateTypeSuccess = global.NotificationTemplateTypeSuccess
	NotificationTemplateTypeFailure = global.NotificationTemplateTypeFailure
)

// NotificationConfig 通知配置
type NotificationConfig struct {
	gorm.Model
	TaskID  uint                `json:"task_id" gorm:"type:int;not null;index:idx_noti_task_id"` // 任务ID，0表示全局配置
	Name    string              `json:"name" gorm:"size:64;not null"`                            // 通知名称
	Type    NotificationType    `json:"type" gorm:"type:varchar(20);not null"`                   // email / webhook
	Target  string              `json:"target" gorm:"type:varchar(255);not null"`                // 邮箱地址或webhook URL
	Trigger NotificationTrigger `json:"trigger" gorm:"type:varchar(20);not null;default:'all'"`  // success / failure / all / timeout / cancel
	Enabled bool                `json:"enabled" gorm:"type:boolean;not null;default:true"`       // 是否启用

	// 邮件配置（type=email时有效）
	EmailRecipients string `json:"email_recipients" gorm:"type:text"` // 收件人列表，逗号分隔（多个邮箱）
	EmailSubject    string `json:"email_subject" gorm:"size:256"`     // 邮件主题模板
	EmailBody       string `json:"email_body" gorm:"size:16777215"`   // 邮件正文模板

	// Webhook配置（type=webhook时有效）
	WebhookURL     string `json:"webhook_url" gorm:"size:512"`                           // Webhook URL
	WebhookMethod  string `json:"webhook_method" gorm:"size:10;not null;default:'POST'"` // HTTP方法
	WebhookHeaders string `json:"webhook_headers" gorm:"type:text"`                      // 自定义请求头 JSON格式
	WebhookBody    string `json:"webhook_body" gorm:"size:16777215"`                     // 自定义请求体模板

	// 重试配置
	RetryTimes    int8 `json:"retry_times" gorm:"not null;default:0"`             // 通知重试次数
	RetryInterval int  `json:"retry_interval" gorm:"type:int;not null;default:5"` // 重试间隔(秒)
}

// NotificationLog 通知发送日志
type NotificationLog struct {
	gorm.Model
	NotificationID uint   `json:"notification_id" gorm:"type:int;not null;index:idx_notification_log_notification_id"` // 通知配置ID
	TaskID         uint   `json:"task_id" gorm:"type:int;not null;index:idx_notification_log_task_id"`                 // 任务ID
	TaskLogID      uint   `json:"task_log_id" gorm:"type:int;not null"`                                                // 任务日志ID
	Status         int8   `json:"status" gorm:"not null;default:1"`                                                    // 发送状态：1=成功, 0=失败
	Result         string `json:"result" gorm:"size:16777215"`                                                         // 发送结果/错误信息
	RetryCount     int8   `json:"retry_count" gorm:"not null;default:0"`                                               // 重试次数
	StartTime      int64  `json:"start_time" gorm:"type:bigint"`                                                       // 开始时间戳
	EndTime        int64  `json:"end_time" gorm:"type:bigint"`                                                         // 结束时间戳
}

// Create 创建通知配置
func (n *NotificationConfig) Create() (uint, error) {
	result := global.GormDB.Create(n)
	if result.Error != nil {
		return 0, result.Error
	}
	return n.ID, nil
}

// Update 更新通知配置
func (n *NotificationConfig) Update(id uint, data g.Map) (int64, error) {
	result := global.GormDB.Model(n).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// GetByTaskId 根据任务ID获取所有已启用的通知配置
func (n *NotificationConfig) GetByTaskId(taskId uint) (configs []NotificationConfig, err error) {
	result := global.GormDB.Where("task_id = ? AND enabled = ?", taskId, true).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// GetByTaskIdAndType 根据任务ID和类型获取通知配置
func (n *NotificationConfig) GetByTaskIdAndType(taskId uint, notificationType NotificationType) (configs []NotificationConfig, err error) {
	result := global.GormDB.Where("task_id = ? AND type = ? AND enabled = ?", taskId, notificationType, true).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// GetGlobalNotifications 获取全局通知配置
func (n *NotificationConfig) GetGlobalNotifications() (configs []NotificationConfig, err error) {
	result := global.GormDB.Where("task_id = 0 AND enabled = ?", true).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// GetNotificationsForTask 获取任务的所有通知配置（包括全局配置）
func (n *NotificationConfig) GetNotificationsForTask(taskId uint) (configs []NotificationConfig, err error) {
	result := global.GormDB.Where("(task_id = ? OR task_id = 0) AND enabled = ?", taskId, true).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// PageList 分页查询通知配置
func (n *NotificationConfig) PageList(tx *gorm.DB, taskId uint, pageNo, pageSize int) (configs []NotificationConfig, count int64, err error) {
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := tx.Model(n)
	if taskId > 0 {
		query = query.Where("task_id = ?", taskId)
	}

	query.Count(&count)

	offset := (pageNo - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("id desc").Find(&configs)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return
}

// FindOne 查找单个通知配置
func (n *NotificationConfig) FindOne(tx *gorm.DB, data g.Map) error {
	result := tx.Where(data).First(n)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID 根据ID查找通知配置
func (n *NotificationConfig) FindByID(id uint) error {
	result := global.GormDB.Where("id = ?", id).First(n)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 删除通知配置
func (n *NotificationConfig) Delete(tx *gorm.DB, id uint) error {
	result := tx.Delete(n, id)
	return result.Error
}

// DeleteByTaskId 删除任务下所有通知配置
func (n *NotificationConfig) DeleteByTaskId(tx *gorm.DB, taskId uint) error {
	result := tx.Where("task_id = ?", taskId).Delete(n)
	return result.Error
}

// GetEnabled 获取所有已启用的通知配置
func (n *NotificationConfig) GetEnabled() (configs []NotificationConfig, err error) {
	result := global.GormDB.Where("enabled = ?", true).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// CreateLog 创建通知发送日志
func (n *NotificationLog) Create() (uint, error) {
	result := global.GormDB.Create(n)
	if result.Error != nil {
		return 0, result.Error
	}
	return n.ID, nil
}

// UpdateLog 更新通知发送日志
func (n *NotificationLog) UpdateLog(id uint, data g.Map) (int64, error) {
	result := global.GormDB.Model(n).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// GetLogsByTaskID 获取任务的通知发送日志
func (n *NotificationLog) GetLogsByTaskID(taskID uint, page, pageSize int) ([]NotificationLog, int64, error) {
	var logs []NotificationLog
	var count int64

	query := global.GormDB.Model(&NotificationLog{}).Where("task_id = ?", taskID)
	query.Count(&count)

	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&logs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return logs, count, nil
}

// GetLogsByNotificationID 获取通知配置的发送日志
func (n *NotificationLog) GetLogsByNotificationID(notificationID uint, page, pageSize int) ([]NotificationLog, int64, error) {
	var logs []NotificationLog
	var count int64

	query := global.GormDB.Model(&NotificationLog{}).Where("notification_id = ?", notificationID)
	query.Count(&count)

	offset := (page - 1) * pageSize
	result := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&logs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return logs, count, nil
}

// ShouldNotify 检查是否应该发送通知
func (n *NotificationConfig) ShouldNotify(trigger NotificationTrigger) bool {
	if !n.Enabled {
		return false
	}

	switch n.Trigger {
	case NotificationTriggerAll:
		return true
	case NotificationTriggerSuccess:
		return trigger == NotificationTriggerSuccess
	case NotificationTriggerFailure:
		return trigger == NotificationTriggerFailure || trigger == NotificationTriggerTimeout || trigger == NotificationTriggerCancel
	case NotificationTriggerTimeout:
		return trigger == NotificationTriggerTimeout
	case NotificationTriggerCancel:
		return trigger == NotificationTriggerCancel
	default:
		return false
	}
}
