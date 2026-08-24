package notification

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"cronJob/internal/global"
	"cronJob/internal/models"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ============================================================
// 邮件/Webhook 通知服务
// ============================================================

// emailConfig 邮件全局配置（从 config.yaml 读取）
type emailConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// webhookConfig Webhook全局配置（从 config.yaml 读取）
type webhookConfig struct {
	Enabled bool              `mapstructure:"enabled"`
	URL     string            `mapstructure:"url"`
	Headers map[string]string `mapstructure:"headers"`
}

// ============================================================
// 通知任务上下文（用于模板渲染）
// ============================================================

// NotifyContext 传递给通知模板的上下文
type NotifyContext struct {
	TaskID     uint
	TaskName   string
	Command    string
	Status     string // "success" / "failure"
	Result     string
	RetryTimes int8
	StartTime  string
	EndTime    string
	Duration   string // 人类可读时长
}

// ============================================================
// 公开方法
// ============================================================

// SendNotification 根据任务执行结果向该任务的所有已启用通知配置发送通知
// 在 afterExecJob 中异步调用，不阻塞主流程
func SendNotification(taskModel *models.Task, taskResult global.TaskResult, startTime time.Time) {
	// 查询该任务的通知配置
	nc := &models.NotificationConfig{}
	configs, err := nc.GetByTaskId(taskModel.ID)
	if err != nil {
		zap.S().Errorf("查询任务通知配置失败: taskId=%d, error=%v", taskModel.ID, err)
		return
	}

	if len(configs) == 0 {
		return
	}

	// 构建通知上下文
	ctx := buildNotifyContext(taskModel, taskResult, startTime)

	// 逐个发送通知（数量通常很少，串行即可）
	for _, cfg := range configs {
		// 检查触发条件是否匹配
		if !shouldTrigger(cfg.Trigger, taskResult.Err != nil) {
			continue
		}

		switch cfg.Type {
		case models.NotificationTypeEmail:
			go sendEmailNotification(cfg, ctx)
		case models.NotificationTypeWebhook:
			go sendWebhookNotification(cfg, ctx)
		default:
			zap.S().Warnf("未知通知类型: %s", cfg.Type)
		}
	}
}

// ============================================================
// 内部实现
// ============================================================

// buildNotifyContext 构建通知上下文
func buildNotifyContext(taskModel *models.Task, taskResult global.TaskResult, startTime time.Time) NotifyContext {
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	status := "success"
	if taskResult.Err != nil {
		status = "failure"
	}

	result := taskResult.Result
	if taskResult.Err != nil {
		result = taskResult.Err.Error()
	}

	// 截断过长的结果，避免邮件内容过大
	const maxResultLen = 2000
	if len(result) > maxResultLen {
		result = result[:maxResultLen] + fmt.Sprintf("...[已截断，共 %d 字节]", len(result))
	}

	return NotifyContext{
		TaskID:     taskModel.ID,
		TaskName:   taskModel.Name,
		Command:    taskModel.Command,
		Status:     status,
		Result:     result,
		RetryTimes: taskResult.RetryTimes,
		StartTime:  startTime.Format("2006-01-02 15:04:05"),
		EndTime:    endTime.Format("2006-01-02 15:04:05"),
		Duration:   formatDuration(duration),
	}
}

// shouldTrigger 判断通知触发条件是否匹配
func shouldTrigger(trigger models.NotificationTrigger, isFailure bool) bool {
	switch trigger {
	case models.NotificationTriggerSuccess:
		return !isFailure
	case models.NotificationTriggerFailure:
		return isFailure
	case models.NotificationTriggerAll:
		return true
	default:
		return true
	}
}

// formatDuration 格式化时长为人类可读格式
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm%ds", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
}

// ============================================================
// 邮件通知
// ============================================================

// emailBodyTemplate 邮件正文 HTML 模板
var emailBodyTemplate = template.Must(template.New("email").Parse(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: 'Microsoft YaHei', Arial, sans-serif; padding: 20px; color: #333;">
    <h2 style="color: {{if eq .Status "success"}}#28a745{{else}}#dc3545{{end}};">
        定时任务通知 - {{if eq .Status "success"}}✅ 执行成功{{else}}❌ 执行失败{{end}}
    </h2>
    <table style="border-collapse: collapse; width: 100%; max-width: 600px;">
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">任务ID</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;">{{.TaskID}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">任务名称</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;">{{.TaskName}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">执行命令</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;"><code>{{.Command}}</code></td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">执行状态</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd; color: {{if eq .Status "success"}}#28a745{{else}}#dc3545{{end}}; font-weight: bold;">
                {{if eq .Status "success"}}✅ 执行成功{{else}}❌ 执行失败{{end}}
            </td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">重试次数</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;">{{.RetryTimes}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">开始时间</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;">{{.StartTime}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">结束时间</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;">{{.EndTime}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; border-bottom: 1px solid #ddd;">执行时长</td>
            <td style="padding: 8px 12px; border-bottom: 1px solid #ddd;">{{.Duration}}</td>
        </tr>
        <tr>
            <td style="padding: 8px 12px; font-weight: bold; vertical-align: top;">执行结果</td>
            <td style="padding: 8px 12px;">
                <pre style="background: #f5f5f5; padding: 10px; border-radius: 4px; white-space: pre-wrap; word-break: break-all; max-height: 300px; overflow-y: auto; font-size: 12px;">{{.Result}}</pre>
            </td>
        </tr>
    </table>
    <p style="color: #999; font-size: 12px; margin-top: 20px;">—— 由 CronJob 定时任务管理系统自动发送</p>
</body>
</html>
`))

// sendEmailNotification 发送邮件通知
func sendEmailNotification(cfg models.NotificationConfig, ctx NotifyContext) {
	// 读取全局邮件配置
	ec := emailConfig{
		Enabled:  viper.GetBool("notification.email.enabled"),
		SMTPHost: viper.GetString("notification.email.smtp_host"),
		SMTPPort: viper.GetInt("notification.email.smtp_port"),
		Username: viper.GetString("notification.email.username"),
		Password: viper.GetString("notification.email.password"),
		From:     viper.GetString("notification.email.from"),
	}

	if !ec.Enabled {
		zap.S().Debugf("邮件通知全局配置未启用，跳过: taskId=%d", ctx.TaskID)
		return
	}

	if ec.SMTPHost == "" || ec.SMTPPort == 0 {
		zap.S().Warnf("邮件SMTP配置不完整，跳过通知: taskId=%d", ctx.TaskID)
		return
	}

	// 渲染邮件内容
	var body bytes.Buffer
	if err := emailBodyTemplate.Execute(&body, ctx); err != nil {
		zap.S().Errorf("渲染邮件模板失败: %v", err)
		return
	}

	// 构建邮件
	subject := fmt.Sprintf("【CronJob】任务「%s」执行%s", ctx.TaskName,
		map[string]string{"success": "成功", "failure": "失败"}[ctx.Status])

	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", ec.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", cfg.Target))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.Write(body.Bytes())

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", ec.SMTPHost, ec.SMTPPort)
	auth := smtp.PlainAuth("", ec.Username, ec.Password, ec.SMTPHost)

	err := sendEmailWithTLS(addr, auth, ec.From, cfg.Target, msg.Bytes())
	if err != nil {
		zap.S().Errorf("发送邮件通知失败: taskId=%d, to=%s, error=%v", ctx.TaskID, cfg.Target, err)
		return
	}

	zap.S().Infof("邮件通知发送成功: taskId=%d, to=%s", ctx.TaskID, cfg.Target)
}

// sendEmailWithTLS 支持 STARTTLS 的邮件发送
func sendEmailWithTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	// 先尝试直接 TLS 连接（端口 465）
	if strings.HasSuffix(addr, ":465") {
		return sendEmailWithDirectTLS(addr, auth, from, to, msg)
	}

	// 尝试普通连接 + STARTTLS（端口 587）
	c, err := smtp.Dial(addr)
	if err != nil {
		// STARTTLS 连接失败，回退到直接 TLS
		return sendEmailWithDirectTLS(addr, auth, from, to, msg)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{
			ServerName: strings.Split(addr, ":")[0],
		}
		if err = c.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("STARTTLS failed: %w", err)
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return fmt.Errorf("SMTP auth failed: %w", err)
			}
		}
	}

	return smtpSendMail(c, from, to, msg)
}

// sendEmailWithDirectTLS 直接TLS连接发送邮件
func sendEmailWithDirectTLS(addr string, auth smtp.Auth, from, to string, msg []byte) error {
	host := strings.Split(addr, ":")[0]
	tlsCfg := &tls.Config{
		ServerName: host,
	}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return fmt.Errorf("SMTP auth failed: %w", err)
			}
		}
	}

	return smtpSendMail(c, from, to, msg)
}

// smtpSendMail 通用的 SMTP 邮件发送
func smtpSendMail(c *smtp.Client, from, to string, msg []byte) error {
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return c.Quit()
}

// ============================================================
// Webhook 通知
// ============================================================

// webhookPayload Webhook 请求体
type webhookPayload struct {
	TaskID     uint   `json:"task_id"`
	TaskName   string `json:"task_name"`
	Command    string `json:"command"`
	Status     string `json:"status"`
	Result     string `json:"result"`
	RetryTimes int8   `json:"retry_times"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Duration   string `json:"duration"`
}

// sendWebhookNotification 发送Webhook通知
func sendWebhookNotification(cfg models.NotificationConfig, ctx NotifyContext) {
	// 读取全局Webhook配置
	wc := webhookConfig{
		Enabled: viper.GetBool("notification.webhook.enabled"),
		URL:     viper.GetString("notification.webhook.url"),
		Headers: viper.GetStringMapString("notification.webhook.headers"),
	}

	// 确定webhook URL：优先使用通知配置中的target，其次使用全局配置
	webhookURL := cfg.Target
	if webhookURL == "" {
		webhookURL = wc.URL
	}

	if webhookURL == "" {
		zap.S().Warnf("Webhook URL为空，跳过通知: taskId=%d", ctx.TaskID)
		return
	}

	// 构建请求体
	payload := webhookPayload{
		TaskID:     ctx.TaskID,
		TaskName:   ctx.TaskName,
		Command:    ctx.Command,
		Status:     ctx.Status,
		Result:     ctx.Result,
		RetryTimes: ctx.RetryTimes,
		StartTime:  ctx.StartTime,
		EndTime:    ctx.EndTime,
		Duration:   ctx.Duration,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		zap.S().Errorf("序列化Webhook请求体失败: %v", err)
		return
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(jsonData))
	if err != nil {
		zap.S().Errorf("创建Webhook请求失败: %v", err)
		return
	}

	// 设置默认Content-Type
	req.Header.Set("Content-Type", "application/json")

	// 设置全局配置中的自定义headers
	for k, v := range wc.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求（超时10秒）
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		zap.S().Errorf("发送Webhook通知失败: taskId=%d, url=%s, error=%v", ctx.TaskID, webhookURL, err)
		return
	}
	defer resp.Body.Close()

	// 读取响应（限制大小）
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		zap.S().Errorf("读取Webhook响应失败: taskId=%d, error=%v", ctx.TaskID, err)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		zap.S().Infof("Webhook通知发送成功: taskId=%d, url=%s, status=%d", ctx.TaskID, webhookURL, resp.StatusCode)
	} else {
		zap.S().Warnf("Webhook通知返回非成功状态: taskId=%d, url=%s, status=%d, body=%s",
			ctx.TaskID, webhookURL, resp.StatusCode, string(respBody))
	}
}
