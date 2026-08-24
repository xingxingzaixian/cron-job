package notification

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"fmt"
	"text/template"
	"bytes"
	"time"
)

// TemplateService 模板渲染服务
type TemplateService struct {
}

// TemplateData 模板数据
type TemplateData struct {
	TaskID    uint
	TaskName  string
	Status    string
	StartTime string
	EndTime   string
	Duration  int64
	Result    string
	Error     string
	Protocol  string
	Command   string
}

// NewTemplateService 创建模板服务实例
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// Render 渲染模板
func (s *TemplateService) Render(templateStr string, task *models.Task, taskLog *models.TaskLog) (string, error) {
	// 准备模板数据
	data := s.prepareData(task, taskLog)

	// 解析模板
	tmpl, err := template.New("notification").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}

	// 执行模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	return buf.String(), nil
}

// prepareData 准备模板数据
func (s *TemplateService) prepareData(task *models.Task, taskLog *models.TaskLog) TemplateData {
	data := TemplateData{
		TaskID:    task.ID,
		TaskName:  task.Name,
		Protocol:  s.getProtocolName(task.Protocol),
		Command:   task.Command,
	}

	if taskLog != nil {
		data.Status = s.getStatusName(taskLog.Status)
		data.StartTime = formatTime(taskLog.StartTime)
		data.EndTime = formatTime(taskLog.EndTime)
		data.Duration = taskLog.Duration

		if taskLog.Status == global.TaskStatusFinish {
			data.Result = taskLog.Result
		} else {
			data.Error = taskLog.Result
		}
	}

	return data
}

// getStatusName 获取状态名称
func (s *TemplateService) getStatusName(status global.TaskStatus) string {
	switch status {
	case global.TaskStatusFinish:
		return "成功"
	case global.TaskStatusFailure:
		return "失败"
	case global.TaskStatusTimeout:
		return "超时"
	case global.TaskStatusCancel:
		return "取消"
	case global.TaskStatusRunning:
		return "运行中"
	case global.TaskStatusEnabled:
		return "启用"
	case global.TaskStatusDisabled:
		return "禁用"
	default:
		return "未知"
	}
}

// getProtocolName 获取协议名称
func (s *TemplateService) getProtocolName(protocol global.TaskProtocol) string {
	switch protocol {
	case global.TaskProtocolHttp:
		return "HTTP"
	case global.TaskProtocolShell:
		return "Shell"
	case global.TaskProtocolSSH:
		return "SSH"
	default:
		return "未知"
	}
}

// RenderString 渲染简单字符串模板
func (s *TemplateService) RenderString(templateStr string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("simple").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	return buf.String(), nil
}

// formatTime 时间格式化
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
