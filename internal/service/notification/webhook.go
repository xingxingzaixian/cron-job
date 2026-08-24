package notification

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// WebhookService Webhook调用服务
type WebhookService struct {
	client      *http.Client
	contentType string
}

// WebhookPayload Webhook请求体
type WebhookPayload struct {
	Event     string `json:"event"`
	TaskID    uint   `json:"task_id"`
	TaskName  string `json:"task_name"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Duration  int64  `json:"duration"`
	Result    string `json:"result,omitempty"`
	Error     string `json:"error,omitempty"`
}

// NewWebhookService 创建Webhook服务实例
func NewWebhookService() *WebhookService {
	timeout := viper.GetInt("notification.webhook.timeout")
	if timeout <= 0 {
		timeout = 10
	}

	return &WebhookService{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		contentType: "application/json",
	}
}

// Send 发送Webhook通知
func (s *WebhookService) Send(notification *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) error {
	// 构建Webhook URL
	webhookURL := s.buildWebhookURL(notification)
	if webhookURL == "" {
		return fmt.Errorf("Webhook URL未配置")
	}

	// 构建请求体
	payload, err := s.buildPayload(notification, task, taskLog)
	if err != nil {
		return fmt.Errorf("构建Webhook请求体失败: %v", err)
	}

	// 发送请求
	return s.sendRequest(webhookURL, notification.WebhookMethod, notification.WebhookHeaders, payload)
}

// buildWebhookURL 构建Webhook URL
func (s *WebhookService) buildWebhookURL(notification *models.NotificationConfig) string {
	// 优先使用通知配置中的URL
	if notification.WebhookURL != "" {
		return notification.WebhookURL
	}

	// 使用Target作为URL
	return notification.Target
}

// buildPayload 构建请求体
func (s *WebhookService) buildPayload(notification *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) (string, error) {
	// 如果有自定义模板，使用模板渲染
	if notification.WebhookBody != "" {
		templateService := NewTemplateService()
		return templateService.Render(notification.WebhookBody, task, taskLog)
	}

	// 使用默认JSON格式
	payload := WebhookPayload{
		Event:     s.getEventName(taskLog.Status),
		TaskID:    task.ID,
		TaskName:  task.Name,
		Status:    s.getStatusName(taskLog.Status),
		StartTime: formatTime(taskLog.StartTime),
		EndTime:   formatTime(taskLog.EndTime),
		Duration:  taskLog.Duration,
	}

	if taskLog.Status == global.TaskStatusFinish {
		payload.Result = taskLog.Result
	} else {
		payload.Error = taskLog.Result
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// getEventName 获取事件名称
func (s *WebhookService) getEventName(status global.TaskStatus) string {
	switch status {
	case global.TaskStatusFinish:
		return "task_success"
	case global.TaskStatusFailure:
		return "task_failure"
	case global.TaskStatusTimeout:
		return "task_timeout"
	case global.TaskStatusCancel:
		return "task_cancel"
	default:
		return "task_unknown"
	}
}

// getStatusName 获取状态名称
func (s *WebhookService) getStatusName(status global.TaskStatus) string {
	switch status {
	case global.TaskStatusFinish:
		return "success"
	case global.TaskStatusFailure:
		return "failure"
	case global.TaskStatusTimeout:
		return "timeout"
	case global.TaskStatusCancel:
		return "cancel"
	default:
		return "unknown"
	}
}

// sendRequest 发送HTTP请求
func (s *WebhookService) sendRequest(url, method, headersStr, payload string) error {
	// 构建请求
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置Content-Type
	contentType := s.contentType
	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)

	// 设置自定义请求头
	if headersStr != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(headersStr), &headers); err == nil {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
	}

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Webhook返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	zap.S().Infof("Webhook发送成功: url=%s, status=%d", url, resp.StatusCode)
	return nil
}

// TestConnection 测试Webhook连接
func (s *WebhookService) TestConnection(url string) error {
	if url == "" {
		return fmt.Errorf("Webhook URL未配置")
	}

	// 发送GET请求测试连接
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("连接Webhook失败: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// TestSend 发送测试Webhook
func (s *WebhookService) TestSend(url, method, headers, payload string) error {
	if url == "" {
		return fmt.Errorf("Webhook URL未配置")
	}

	if method == "" {
		method = "POST"
	}

	return s.sendRequest(url, method, headers, payload)
}

// 确保 WebhookService 实现了 Sender 接口
var _ Sender = (*WebhookService)(nil)
