package notification

import (
	"cronJob/internal/models"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NotificationManager 通知管理器
type NotificationManager struct {
	// 缓存的通知配置
	cache      map[uint][]models.NotificationConfig
	cacheMutex sync.RWMutex
	cacheTime  time.Time
	cacheTTL   time.Duration
}

// NewNotificationManager 创建通知管理器实例
func NewNotificationManager() *NotificationManager {
	cacheTTL := viper.GetInt("notification.cache.ttl")
	if cacheTTL <= 0 {
		cacheTTL = 60 // 默认60秒
	}

	return &NotificationManager{
		cache:     make(map[uint][]models.NotificationConfig),
		cacheTTL:  time.Duration(cacheTTL) * time.Second,
		cacheTime: time.Now(),
	}
}

// GetNotificationsForTask 获取任务的所有通知配置（包括全局配置）
func (m *NotificationManager) GetNotificationsForTask(taskID uint) ([]models.NotificationConfig, error) {
	// 检查缓存
	m.cacheMutex.RLock()
	if time.Since(m.cacheTime) < m.cacheTTL {
		if configs, ok := m.cache[taskID]; ok {
			m.cacheMutex.RUnlock()
			return configs, nil
		}
	}
	m.cacheMutex.RUnlock()

	// 从数据库获取
	notification := &models.NotificationConfig{}
	configs, err := notification.GetNotificationsForTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("获取任务通知配置失败: %v", err)
	}

	// 更新缓存
	m.cacheMutex.Lock()
	m.cache[taskID] = configs
	m.cacheTime = time.Now()
	m.cacheMutex.Unlock()

	return configs, nil
}

// GetNotificationsByType 获取指定类型的通知配置
func (m *NotificationManager) GetNotificationsByType(taskID uint, notificationType models.NotificationType) ([]models.NotificationConfig, error) {
	// 获取所有配置
	configs, err := m.GetNotificationsForTask(taskID)
	if err != nil {
		return nil, err
	}

	// 过滤指定类型
	var filtered []models.NotificationConfig
	for _, config := range configs {
		if config.Type == notificationType && config.Enabled {
			filtered = append(filtered, config)
		}
	}

	return filtered, nil
}

// ShouldNotify 检查是否应该发送通知
func (m *NotificationManager) ShouldNotify(config models.NotificationConfig, trigger models.NotificationTrigger) bool {
	return config.ShouldNotify(trigger)
}

// ClearCache 清除缓存
func (m *NotificationManager) ClearCache() {
	m.cacheMutex.Lock()
	m.cache = make(map[uint][]models.NotificationConfig)
	m.cacheTime = time.Time{}
	m.cacheMutex.Unlock()
}

// ClearCacheByTaskID 清除指定任务的缓存
func (m *NotificationManager) ClearCacheByTaskID(taskID uint) {
	m.cacheMutex.Lock()
	delete(m.cache, taskID)
	m.cacheMutex.Unlock()
}

// ReloadConfig 重新加载配置
func (m *NotificationManager) ReloadConfig() error {
	m.ClearCache()
	zap.S().Info("通知配置缓存已清除")
	return nil
}

// ValidateNotificationConfig 验证通知配置
func (m *NotificationManager) ValidateNotificationConfig(config *models.NotificationConfig) error {
	// 验证通知类型
	switch config.Type {
	case models.NotificationTypeEmail:
		// 验证邮件配置
		if config.EmailRecipients == "" && config.Target == "" {
			return fmt.Errorf("邮件收件人不能为空")
		}
	case models.NotificationTypeWebhook:
		// 验证Webhook配置
		if config.WebhookURL == "" && config.Target == "" {
			return fmt.Errorf("Webhook URL不能为空")
		}
	default:
		return fmt.Errorf("不支持的通知类型: %s", config.Type)
	}

	// 验证触发条件（如果为空则设置默认值）
	if config.Trigger == "" {
		config.Trigger = models.NotificationTriggerAll
	}
	
	switch config.Trigger {
	case models.NotificationTriggerSuccess, models.NotificationTriggerFailure, 
		models.NotificationTriggerAll, models.NotificationTriggerTimeout, models.NotificationTriggerCancel:
		// 有效的触发条件
	default:
		return fmt.Errorf("不支持的触发条件: %s", config.Trigger)
	}

	// 验证重试配置
	if config.RetryTimes < 0 {
		return fmt.Errorf("重试次数不能为负数")
	}
	if config.RetryInterval < 0 {
		return fmt.Errorf("重试间隔不能为负数")
	}

	return nil
}

// GetDefaultNotificationConfig 获取默认通知配置
func (m *NotificationManager) GetDefaultNotificationConfig() models.NotificationConfig {
	return models.NotificationConfig{
		Type:          models.NotificationTypeEmail,
		Trigger:       models.NotificationTriggerAll,
		Enabled:       true,
		RetryTimes:    0,
		RetryInterval: 5,
		WebhookMethod: "POST",
	}
}

// TestNotification 测试通知发送
func (m *NotificationManager) TestNotification(config *models.NotificationConfig) error {
	// 验证配置
	if err := m.ValidateNotificationConfig(config); err != nil {
		return err
	}

	// 根据类型发送测试通知
	switch config.Type {
	case models.NotificationTypeEmail:
		emailService := NewEmailService()
		return emailService.TestConnection()
	case models.NotificationTypeWebhook:
		webhookService := NewWebhookService()
		url := config.WebhookURL
		if url == "" {
			url = config.Target
		}
		return webhookService.TestConnection(url)
	default:
		return fmt.Errorf("不支持的通知类型: %s", config.Type)
	}
}
