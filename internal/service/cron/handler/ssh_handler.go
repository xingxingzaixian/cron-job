package handler

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/service/cron/lib/sshclient"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SSHHandler SSH远程执行任务
type SSHHandler struct{}

// SSHConfig SSH连接配置结构（存储在Params字段中）
type SSHConfig struct {
	Host     string `json:"host"`     // 主机地址
	Port     int    `json:"port"`     // 端口，默认22
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码（密码认证时使用）
	Mode     string `json:"mode"`     // 执行模式：sequential(顺序执行) 或 script(脚本执行)
}

// Run 实现Handler接口，执行SSH任务
func (h *SSHHandler) Run(taskModel *models.Task, taskUniqueId uint) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error("SSH Handler panic: ", err)
		}
	}()

	// 设置默认超时时间
	if taskModel.Timeout <= 0 {
		taskModel.Timeout = global.SSHExecTimeout
	}

	// 解析SSH配置（从Params字段）
	sshConfig, err := h.parseSSHConfig(taskModel)
	if err != nil {
		return "", fmt.Errorf("任务【%d】SSH配置解析失败: %v", taskModel.ID, err)
	}

	// 解析要执行的命令（从Command字段）
	commands, err := h.parseCommands(taskModel.Command)
	if err != nil {
		return "", fmt.Errorf("任务【%d】命令解析失败: %v", taskModel.ID, err)
	}

	// 验证SSH配置和命令
	if err := h.validateSSHConfig(sshConfig, commands); err != nil {
		return "", fmt.Errorf("任务【%d】SSH配置验证失败: %v", taskModel.ID, err)
	}

	zap.S().Infof("execute ssh start: [id: %d host: %s:%d user: %s commands: %d]",
		taskUniqueId, sshConfig.Host, sshConfig.Port, sshConfig.Username, len(commands))

	// 创建SSH客户端
	client := sshclient.NewSSHClient(&sshclient.SSHConfig{
		Host:     sshConfig.Host,
		Port:     sshConfig.Port,
		Username: sshConfig.Username,
		Password: sshConfig.Password,
		Timeout:  global.SSHConnTimeout,
	})

	// 建立连接
	if err := client.Connect(); err != nil {
		return "", fmt.Errorf("任务【%d】SSH连接失败: %v", taskModel.ID, err)
	}
	defer client.Close()

	// 创建执行上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(taskModel.Timeout)*time.Second)
	defer cancel()

	// 根据执行模式执行命令
	var output string
	switch strings.ToLower(sshConfig.Mode) {
	case "script":
		// 脚本模式：将所有命令作为一个脚本执行
		output, err = h.executeAsScript(ctx, client, commands)
	case "sequential", "":
		// 顺序模式：逐条执行命令（默认模式）
		output, err = h.executeSequentially(ctx, client, commands)
	default:
		return "", fmt.Errorf("任务【%d】不支持的执行模式: %s", taskModel.ID, sshConfig.Mode)
	}

	if err != nil {
		return "", fmt.Errorf("任务【%d】SSH命令执行失败: %v", taskModel.ID, err)
	}

	zap.S().Infof("execute ssh finish: [id: %d host: %s:%d user: %s]",
		taskUniqueId, sshConfig.Host, sshConfig.Port, sshConfig.Username)

	return output, nil
}

// parseSSHConfig 解析SSH配置（从Params字段）
func (h *SSHHandler) parseSSHConfig(taskModel *models.Task) (*SSHConfig, error) {
	config := &SSHConfig{}

	// 解析Params字段中的JSON配置
	if taskModel.Params == "" {
		return nil, fmt.Errorf("SSH连接配置不能为空")
	}

	if err := json.Unmarshal([]byte(taskModel.Params), config); err != nil {
		return nil, fmt.Errorf("SSH配置JSON解析失败: %v", err)
	}

	// 设置默认值
	if config.Port == 0 {
		config.Port = 22
	}
	if config.Mode == "" {
		config.Mode = "sequential"
	}

	return config, nil
}

// parseCommands 解析要执行的命令（从Command字段）
func (h *SSHHandler) parseCommands(commandText string) ([]string, error) {
	if strings.TrimSpace(commandText) == "" {
		return nil, fmt.Errorf("命令内容不能为空")
	}

	// 按行分割命令，支持多行命令
	lines := strings.Split(commandText, "\n")
	var commands []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") { // 忽略空行和注释行
			commands = append(commands, line)
		}
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("没有找到有效的命令")
	}

	return commands, nil
}

// validateSSHConfig 验证SSH配置和命令
func (h *SSHHandler) validateSSHConfig(config *SSHConfig, commands []string) error {
	if config.Host == "" {
		return fmt.Errorf("主机地址不能为空")
	}
	if config.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if config.Password == "" {
		return fmt.Errorf("必须提供密码进行认证")
	}
	if len(commands) == 0 {
		return fmt.Errorf("命令列表不能为空")
	}
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("端口号必须在1-65535之间")
	}

	return nil
}

// executeSequentially 顺序执行命令（在同一个session中保持上下文）
func (h *SSHHandler) executeSequentially(ctx context.Context, client *sshclient.SSHClient, commands []string) (string, error) {
	zap.S().Debugf("执行SSH命令序列，共 %d 条命令，保持session上下文", len(commands))

	// 使用ExecuteCommands在同一个session中执行所有命令，保持上下文
	outputs, err := client.ExecuteCommands(ctx, commands)
	if err != nil {
		return "", fmt.Errorf("命令序列执行失败: %v", err)
	}

	// 格式化输出结果
	var results []string
	for i, output := range outputs {
		result := fmt.Sprintf("~# %d: %s \n%s", i+1, commands[i], strings.TrimSpace(output))
		results = append(results, result)
	}

	return strings.Join(results, "\n\n"), nil
}

// executeAsScript 将命令作为脚本执行
func (h *SSHHandler) executeAsScript(ctx context.Context, client *sshclient.SSHClient, commands []string) (string, error) {
	zap.S().Debugf("执行SSH脚本，共 %d 条命令", len(commands))

	output, err := client.ExecuteScript(ctx, commands)
	if err != nil {
		return "", fmt.Errorf("脚本执行失败: %v", err)
	}

	// 格式化输出
	result := fmt.Sprintf("=== 脚本执行结果 ===\n执行的命令:\n%s\n\n输出:\n%s",
		strings.Join(commands, "\n"), strings.TrimSpace(output))

	return result, nil
}

// TestSSHConnection 测试SSH连接（可用于配置验证）
func TestSSHConnection(config *SSHConfig, commands []string) error {
	if err := (&SSHHandler{}).validateSSHConfig(config, commands); err != nil {
		return err
	}

	client := sshclient.NewSSHClient(&sshclient.SSHConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Timeout:  global.SSHConnTimeout,
	})

	if err := client.Connect(); err != nil {
		return fmt.Errorf("SSH连接测试失败: %v", err)
	}
	defer client.Close()

	// 执行一个简单的测试命令
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.ExecuteCommand(ctx, "echo 'SSH连接测试成功'")
	return err
}

// GetSSHConfigExample 获取SSH配置示例
func GetSSHConfigExample() (string, string) {
	// SSH连接配置（存储在Params字段）
	configExample := SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "your_password", // 或者使用私钥认证
		Mode:     "sequential",    // 或 "script"
	}

	configData, _ := json.MarshalIndent(configExample, "", "  ")

	// 命令示例（存储在Command字段）
	commandExample := `whoami
pwd
ls -la
df -h
# 这是注释，会被忽略
uptime`

	return string(configData), commandExample
}
