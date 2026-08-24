package notification

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// EmailService 邮件发送服务
type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromAddress  string
	fromName     string
	useTLS       bool
	timeout      int
}

// NewEmailService 创建邮件服务实例
func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     viper.GetString("notification.email.smtp_host"),
		smtpPort:     viper.GetInt("notification.email.smtp_port"),
		smtpUsername: viper.GetString("notification.email.username"),
		smtpPassword: viper.GetString("notification.email.password"),
		fromAddress:  viper.GetString("notification.email.from"),
		fromName:     viper.GetString("notification.email.from"),
		useTLS:       true,
		timeout:      30,
	}
}

// Send 发送邮件通知
func (s *EmailService) Send(notification *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) error {
	// 检查邮件服务是否启用
	if !viper.GetBool("notification.email.enabled") {
		return fmt.Errorf("邮件服务未启用")
	}

	// 验证配置
	if err := s.validateConfig(); err != nil {
		return fmt.Errorf("邮件配置无效: %v", err)
	}

	// 构建收件人列表
	recipients := s.buildRecipients(notification)
	if len(recipients) == 0 {
		return fmt.Errorf("没有有效的收件人")
	}

	// 渲染邮件内容
	subject, body, err := s.renderEmailContent(notification, task, taskLog)
	if err != nil {
		return fmt.Errorf("渲染邮件内容失败: %v", err)
	}

	// 发送邮件
	return s.sendEmail(recipients, subject, body)
}

// validateConfig 验证邮件配置
func (s *EmailService) validateConfig() error {
	if s.smtpHost == "" {
		return fmt.Errorf("SMTP主机未配置")
	}
	if s.smtpPort <= 0 {
		return fmt.Errorf("SMTP端口未配置")
	}
	if s.smtpUsername == "" {
		return fmt.Errorf("SMTP用户名未配置")
	}
	if s.smtpPassword == "" {
		return fmt.Errorf("SMTP密码未配置")
	}
	if s.fromAddress == "" {
		return fmt.Errorf("发件人地址未配置")
	}
	return nil
}

// buildRecipients 构建收件人列表
func (s *EmailService) buildRecipients(notification *models.NotificationConfig) []string {
	var recipients []string

	// 从通知配置中获取收件人
	if notification.EmailRecipients != "" {
		parts := strings.Split(notification.EmailRecipients, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				recipients = append(recipients, part)
			}
		}
	}

	// 如果通知配置中有Target且是邮箱格式，也添加到收件人列表
	if notification.Target != "" && strings.Contains(notification.Target, "@") {
		// 检查是否已经在收件人列表中
		found := false
		for _, r := range recipients {
			if r == notification.Target {
				found = true
				break
			}
		}
		if !found {
			recipients = append(recipients, notification.Target)
		}
	}

	return recipients
}

// renderEmailContent 渲染邮件内容
func (s *EmailService) renderEmailContent(notification *models.NotificationConfig, task *models.Task, taskLog *models.TaskLog) (subject string, body string, err error) {
	// 使用模板服务渲染内容
	templateService := NewTemplateService()
	
	// 确定模板类型
	var templateType models.NotificationTemplateType
	if taskLog.Status == global.TaskStatusFinish {
		templateType = models.NotificationTemplateTypeSuccess
	} else {
		templateType = models.NotificationTemplateTypeFailure
	}

	// 渲染主题
	subjectTemplate := notification.EmailSubject
	if subjectTemplate == "" {
		// 使用默认主题
		if templateType == models.NotificationTemplateTypeSuccess {
			subjectTemplate = "任务执行成功: {{.TaskName}}"
		} else {
			subjectTemplate = "任务执行失败: {{.TaskName}}"
		}
	}

	subject, err = templateService.Render(subjectTemplate, task, taskLog)
	if err != nil {
		return "", "", fmt.Errorf("渲染邮件主题失败: %v", err)
	}

	// 渲染正文
	bodyTemplate := notification.EmailBody
	if bodyTemplate == "" {
		// 使用默认正文模板
		bodyTemplate = s.getDefaultBodyTemplate(templateType)
	}

	body, err = templateService.Render(bodyTemplate, task, taskLog)
	if err != nil {
		return "", "", fmt.Errorf("渲染邮件正文失败: %v", err)
	}

	return subject, body, nil
}

// getDefaultBodyTemplate 获取默认邮件正文模板
func (s *EmailService) getDefaultBodyTemplate(templateType models.NotificationTemplateType) string {
	if templateType == models.NotificationTemplateTypeSuccess {
		return `<h2>任务执行成功通知</h2>
<p><strong>任务名称:</strong> {{.TaskName}}</p>
<p><strong>任务ID:</strong> {{.TaskID}}</p>
<p><strong>执行状态:</strong> 成功</p>
<p><strong>开始时间:</strong> {{.StartTime}}</p>
<p><strong>结束时间:</strong> {{.EndTime}}</p>
<p><strong>执行时长:</strong> {{.Duration}}ms</p>
{{if .Result}}<p><strong>执行结果:</strong> <pre>{{.Result}}</pre></p>{{end}}`
	}
	
	return `<h2>任务执行失败通知</h2>
<p><strong>任务名称:</strong> {{.TaskName}}</p>
<p><strong>任务ID:</strong> {{.TaskID}}</p>
<p><strong>执行状态:</strong> 失败</p>
<p><strong>开始时间:</strong> {{.StartTime}}</p>
<p><strong>结束时间:</strong> {{.EndTime}}</p>
<p><strong>执行时长:</strong> {{.Duration}}ms</p>
{{if .Result}}<p><strong>执行结果:</strong> <pre>{{.Result}}</pre></p>{{end}}
{{if .Error}}<p><strong>错误信息:</strong> <pre>{{.Error}}</pre></p>{{end}}`
}

// sendEmail 发送邮件
func (s *EmailService) sendEmail(recipients []string, subject, body string) error {
	// 构建邮件头
	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress),
		"To":           strings.Join(recipients, ","),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	// 构建邮件内容
	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	// 构建SMTP认证
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// 构建服务器地址
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	// 发送邮件
	err := smtp.SendMail(addr, auth, s.fromAddress, recipients, []byte(message.String()))
	if err != nil {
		zap.S().Errorf("发送邮件失败: %v", err)
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	zap.S().Infof("邮件发送成功: recipients=%v, subject=%s", recipients, subject)
	return nil
}

// TestConnection 测试邮件连接
func (s *EmailService) TestConnection() error {
	if err := s.validateConfig(); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	zap.S().Infof("测试邮件连接: %s", addr)

	// 根据端口选择连接方式
	if s.smtpPort == 465 {
		// SSL/TLS 端口，直接使用 TLS 连接
		return s.testTLSConnection(addr)
	} else if s.smtpPort == 587 {
		// STARTTLS 端口
		return s.testSTARTTLSConnection(addr)
	}

	// 其他端口，尝试普通连接
	return s.testPlainConnection(addr)
}

// testTLSConnection 测试 TLS 连接（端口 465）
func (s *EmailService) testTLSConnection(addr string) error {
	host := s.smtpHost
	tlsCfg := &tls.Config{
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return fmt.Errorf("TLS连接失败: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 尝试认证
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	zap.S().Infof("邮件连接测试成功 (TLS)")
	return nil
}

// testSTARTTLSConnection 测试 STARTTLS 连接（端口 587）
func (s *EmailService) testSTARTTLSConnection(addr string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer client.Close()

	// 升级到 TLS
	host := s.smtpHost
	tlsCfg := &tls.Config{
		ServerName: host,
	}
	if err = client.StartTLS(tlsCfg); err != nil {
		return fmt.Errorf("STARTTLS升级失败: %v", err)
	}

	// 尝试认证
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	zap.S().Infof("邮件连接测试成功 (STARTTLS)")
	return nil
}

// testPlainConnection 测试普通连接
func (s *EmailService) testPlainConnection(addr string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer client.Close()

	// 尝试认证
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	zap.S().Infof("邮件连接测试成功 (普通连接)")
	return nil
}

// TestSend 发送测试邮件
func (s *EmailService) TestSend(to, subject, body string) error {
	if err := s.validateConfig(); err != nil {
		return err
	}

	recipients := []string{to}
	return s.sendEmail(recipients, subject, body)
}

// 确保 EmailService 实现了 Sender 接口
var _ Sender = (*EmailService)(nil)
