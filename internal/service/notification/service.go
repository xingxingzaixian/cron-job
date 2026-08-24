package notification

import (
	"cronJob/internal/models"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Sender 通知发送接口
type Sender interface {
	Send(notification *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) error
}

// NotificationService 通知服务
type NotificationService struct {
	emailService    *EmailService
	webhookService  *WebhookService
	templateService *TemplateService
	manager         *NotificationManager
}

var (
	notificationService *NotificationService
	once                sync.Once
)

// GetNotificationService 获取通知服务单例
func GetNotificationService() *NotificationService {
	once.Do(func() {
		notificationService = &NotificationService{
			emailService:    NewEmailService(),
			webhookService:  NewWebhookService(),
			templateService: NewTemplateService(),
			manager:         NewNotificationManager(),
		}
	})
	return notificationService
}

// Notify 发送通知
func (s *NotificationService) Notify(task *models.Task, taskLog *models.TaskLog, trigger models.NotificationTrigger) {
	// 检查通知是否启用
	if !viper.GetBool("notification.email.enabled") && !viper.GetBool("notification.webhook.enabled") {
		return
	}

	// 异步发送通知，不阻塞任务执行
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.S().Errorf("通知发送发生panic: taskId=%d, error=%v", task.ID, r)
			}
		}()

		s.sendNotifications(task, taskLog, trigger)
	}()
}

// sendNotifications 发送所有通知
func (s *NotificationService) sendNotifications(task *models.Task, taskLog *models.TaskLog, trigger models.NotificationTrigger) {
	// 获取任务的通知配置
	configs, err := s.manager.GetNotificationsForTask(task.ID)
	if err != nil {
		zap.S().Errorf("获取任务通知配置失败: taskId=%d, error=%v", task.ID, err)
		return
	}

	// 遍历通知配置
	for _, config := range configs {
		// 检查是否应该发送通知
		if !s.manager.ShouldNotify(config, trigger) {
			continue
		}

		// 发送通知
		go s.sendSingleNotification(&config, task, taskLog)
	}
}

// sendSingleNotification 发送单个通知
func (s *NotificationService) sendSingleNotification(config *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) {
	startTime := time.Now().Unix()

	var err error
	var retryCount int8

	for {
		// 根据类型选择发送器
		var sender Sender
		switch config.Type {
		case models.NotificationTypeEmail:
			sender = s.emailService
		case models.NotificationTypeWebhook:
			sender = s.webhookService
		default:
			zap.S().Errorf("不支持的通知类型: %s", config.Type)
			return
		}

		// 发送通知
		err = sender.Send(config, task, taskLog)
		if err == nil {
			// 发送成功
			s.logNotificationSuccess(config.ID, task.ID, taskLog.ID, retryCount, startTime)
			return
		}

		retryCount++
		if retryCount > config.RetryTimes {
			// 重试次数用完
			s.logNotificationFailure(config.ID, task.ID, taskLog.ID, err.Error(), retryCount, startTime)
			zap.S().Errorf("通知发送失败(已重试%d次): notificationId=%d, taskId=%d, error=%v", 
				config.RetryTimes, config.ID, task.ID, err)
			return
		}

		zap.S().Warnf("通知发送失败，准备重试: notificationId=%d, retryCount=%d, error=%v", 
			config.ID, retryCount, err)
		
		// 等待重试间隔
		time.Sleep(time.Duration(config.RetryInterval) * time.Second)
	}
}

// logNotificationSuccess 记录通知发送成功日志
func (s *NotificationService) logNotificationSuccess(notificationID, taskID, taskLogID uint, retryCount int8, startTime int64) {
	log := &models.NotificationLog{
		NotificationID: notificationID,
		TaskID:         taskID,
		TaskLogID:      taskLogID,
		Status:         1, // 成功
		Result:         "发送成功",
		RetryCount:     retryCount,
		StartTime:      startTime,
		EndTime:        time.Now().Unix(),
	}

	if _, err := log.Create(); err != nil {
		zap.S().Errorf("记录通知发送日志失败: notificationId=%d, error=%v", notificationID, err)
	}
}

// logNotificationFailure 记录通知发送失败日志
func (s *NotificationService) logNotificationFailure(notificationID, taskID, taskLogID uint, result string, retryCount int8, startTime int64) {
	log := &models.NotificationLog{
		NotificationID: notificationID,
		TaskID:         taskID,
		TaskLogID:      taskLogID,
		Status:         0, // 失败
		Result:         result,
		RetryCount:     retryCount,
		StartTime:      startTime,
		EndTime:        time.Now().Unix(),
	}

	if _, err := log.Create(); err != nil {
		zap.S().Errorf("记录通知发送日志失败: notificationId=%d, error=%v", notificationID, err)
	}
}

// GetManager 获取通知管理器
func (s *NotificationService) GetManager() *NotificationManager {
	return s.manager
}

// TestNotification 测试通知发送
func (s *NotificationService) TestNotification(config *models.NotificationConfig) error {
	return s.manager.TestNotification(config)
}

// SendTestNotification 发送测试通知
func (s *NotificationService) SendTestNotification(config *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) error {
	// 验证配置
	if err := s.manager.ValidateNotificationConfig(config); err != nil {
		return err
	}

	// 发送通知
	s.sendSingleNotification(config, task, taskLog)
	return nil
}

// GetNotificationStats 获取通知统计信息
func (s *NotificationService) GetNotificationStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取通知配置数量
	notification := &models.NotificationConfig{}
	configs, err := notification.GetEnabled()
	if err != nil {
		return nil, fmt.Errorf("获取通知配置失败: %v", err)
	}

	stats["total_configs"] = len(configs)

	// 按类型统计
	emailCount := 0
	webhookCount := 0
	for _, config := range configs {
		switch config.Type {
		case models.NotificationTypeEmail:
			emailCount++
		case models.NotificationTypeWebhook:
			webhookCount++
		}
	}

	stats["email_configs"] = emailCount
	stats["webhook_configs"] = webhookCount

	return stats, nil
}
