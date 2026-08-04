package sshclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// SSHConfig SSH连接配置
type SSHConfig struct {
	Host     string // 主机地址
	Port     int    // 端口，默认22
	Username string // 用户名
	Password string // 密码（密码认证时使用）
	Timeout  int    // 连接超时时间（秒）
	HostKey  string // 主机公钥（可选，如 "ssh-rsa AAAA..."），提供时校验主机身份
}

// SSHClient SSH客户端
type SSHClient struct {
	config *SSHConfig
	client *ssh.Client
}

// NewSSHClient 创建SSH客户端
func NewSSHClient(config *SSHConfig) *SSHClient {
	if config.Port == 0 {
		config.Port = 22
	}
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	return &SSHClient{
		config: config,
	}
}

// Connect 建立SSH连接
func (c *SSHClient) Connect() error {
	// 创建SSH客户端配置
	sshConfig := &ssh.ClientConfig{
		User:    c.config.Username,
		Timeout: time.Duration(c.config.Timeout) * time.Second,
	}

	// 配置主机密钥校验：提供了主机公钥则严格校验，否则降级为跳过校验（记录警告）
	if c.config.HostKey != "" {
		hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(c.config.HostKey))
		if err != nil {
			return fmt.Errorf("主机公钥格式错误: %v", err)
		}
		sshConfig.HostKeyCallback = ssh.FixedHostKey(hostKey)
	} else {
		zap.S().Warnf("SSH连接[%s:%d]未配置主机公钥，跳过主机身份校验，存在中间人攻击风险", c.config.Host, c.config.Port)
		sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// 根据认证方式设置认证方法
	if c.config.Password != "" {
		// 密码认证
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.Password(c.config.Password),
		}
	} else {
		return fmt.Errorf("必须提供密码进行认证")
	}

	// 建立连接
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("SSH连接失败: %v", err)
	}

	c.client = client
	return nil
}

// ExecuteCommand 执行单条命令
func (c *SSHClient) ExecuteCommand(ctx context.Context, command string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("SSH连接未建立")
	}

	// 创建会话
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 使用通道来处理超时
	type result struct {
		output string
		err    error
	}

	resultChan := make(chan result, 1)
	go func() {
		output, err := session.CombinedOutput(command)
		resultChan <- result{string(output), err}
	}()

	// 等待命令执行完成或超时
	select {
	case <-ctx.Done():
		// 超时或取消
		session.Signal(ssh.SIGKILL)
		return "", fmt.Errorf("命令执行超时或被取消")
	case res := <-resultChan:
		return res.output, res.err
	}
}

// ExecuteCommands 在同一个session中执行多条命令（保持上下文）
func (c *SSHClient) ExecuteCommands(ctx context.Context, commands []string) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("SSH连接未建立")
	}

	// 创建一个session来执行所有命令
	session, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 构建一个脚本，每个命令后添加分隔符
	var scriptLines []string
	for i, command := range commands {
		delimiter := fmt.Sprintf("===CMD_%d_OUTPUT_END===", i)
		scriptLines = append(scriptLines, command)
		scriptLines = append(scriptLines, fmt.Sprintf("echo '%s'", delimiter))
	}

	script := strings.Join(scriptLines, "\n")

	// 使用通道来处理超时
	type result struct {
		outputs []string
		err     error
	}

	resultChan := make(chan result, 1)
	go func() {
		output, err := session.CombinedOutput(script)
		if err != nil {
			resultChan <- result{nil, fmt.Errorf("脚本执行失败: %v", err)}
			return
		}

		// 按分隔符分割输出
		outputs := c.parseCommandOutputs(string(output), len(commands))
		resultChan <- result{outputs, nil}
	}()

	// 等待命令执行完成或超时
	select {
	case <-ctx.Done():
		session.Signal(ssh.SIGKILL)
		return nil, fmt.Errorf("命令执行超时或被取消")
	case res := <-resultChan:
		return res.outputs, res.err
	}
}

// parseCommandOutputs 解析命令输出，按分隔符分割
func (c *SSHClient) parseCommandOutputs(allOutput string, commandCount int) []string {
	outputs := make([]string, commandCount)

	// 按分隔符分割输出
	for i := 0; i < commandCount; i++ {
		delimiter := fmt.Sprintf("===CMD_%d_OUTPUT_END===", i)

		// 查找当前命令的分隔符
		delimiterIndex := strings.Index(allOutput, delimiter)
		if delimiterIndex == -1 {
			// 如果找不到分隔符，将剩余内容作为最后一个命令的输出
			outputs[i] = strings.TrimSpace(allOutput)
			break
		}

		// 提取当前命令的输出
		commandOutput := allOutput[:delimiterIndex]
		outputs[i] = strings.TrimSpace(commandOutput)

		// 移动到下一个命令的输出
		allOutput = allOutput[delimiterIndex+len(delimiter):]

		// 移除分隔符后的换行符
		if strings.HasPrefix(allOutput, "\n") {
			allOutput = allOutput[1:]
		}
	}

	return outputs
}

// ExecuteScript 执行脚本（多条命令作为一个脚本执行）
func (c *SSHClient) ExecuteScript(ctx context.Context, commands []string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("SSH连接未建立")
	}

	// 将多条命令合并为一个脚本
	script := strings.Join(commands, "\n")

	// 创建会话
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建SSH会话失败: %v", err)
	}
	defer session.Close()

	// 使用通道来处理超时
	type result struct {
		output string
		err    error
	}

	resultChan := make(chan result, 1)
	go func() {
		output, err := session.CombinedOutput(script)
		resultChan <- result{string(output), err}
	}()

	// 等待脚本执行完成或超时
	select {
	case <-ctx.Done():
		// 超时或取消
		session.Signal(ssh.SIGKILL)
		return "", fmt.Errorf("脚本执行超时或被取消")
	case res := <-resultChan:
		return res.output, res.err
	}
}

// Close 关闭SSH连接
func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// IsConnected 检查连接状态
func (c *SSHClient) IsConnected() bool {
	if c.client == nil {
		return false
	}

	// 通过发送一个简单的请求来检查连接状态
	_, _, err := c.client.SendRequest("keepalive", false, nil)
	return err == nil
}

// TestConnection 测试SSH连接
func TestConnection(config *SSHConfig) error {
	client := NewSSHClient(config)
	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Close()

	// 执行一个简单的命令来测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.ExecuteCommand(ctx, "echo 'connection test'")
	return err
}

// GetHostInfo 获取远程主机信息
func (c *SSHClient) GetHostInfo(ctx context.Context) (map[string]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("SSH连接未建立")
	}

	info := make(map[string]string)

	// 获取系统信息的命令
	commands := map[string]string{
		"hostname": "hostname",
		"os":       "uname -s",
		"kernel":   "uname -r",
		"arch":     "uname -m",
		"uptime":   "uptime",
	}

	for key, cmd := range commands {
		output, err := c.ExecuteCommand(ctx, cmd)
		if err != nil {
			info[key] = fmt.Sprintf("获取失败: %v", err)
		} else {
			info[key] = strings.TrimSpace(output)
		}
	}

	return info, nil
}
